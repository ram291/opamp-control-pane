package supervisor

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/open-telemetry/opamp-go/client"
	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
)

// Supervisor manages the OpenTelemetry Collector lifecycle and OpAMP communication.
type Supervisor struct {
	logger       types.Logger
	config       *Config
	opampClient  client.OpAMPClient
	upgrader     *BinaryUpgrader
	cmd          *exec.Cmd
	mu           sync.Mutex
	running      atomic.Bool
	instanceId   uuid.UUID
	agentVersion string
	startedAt    time.Time
}

// New creates a new Supervisor.
func New(logger types.Logger, cfg *Config) *Supervisor {
	return &Supervisor{
		logger:       logger,
		config:       cfg,
		upgrader:     NewBinaryUpgrader(cfg.Artifactory),
		instanceId:   uuid.New(),
		agentVersion: "1.0.0",
	}
}

// Start starts the supervisor: connects to OpAMP server and manages the collector.
func (s *Supervisor) Start(ctx context.Context) error {
	s.logger.Debugf(ctx, "Supervisor starting, id=%s", s.instanceId.String())

	s.opampClient = client.NewWebSocket(s.logger)

	callbacks := types.Callbacks{
		OnConnect: func(ctx context.Context) {
			s.logger.Debugf(ctx, "Connected to OpAMP server.")
		},
		OnConnectFailed: func(ctx context.Context, err error) {
			s.logger.Errorf(ctx, "Failed to connect to OpAMP server: %v", err)
		},
		OnError: func(ctx context.Context, err *protobufs.ServerErrorResponse) {
			s.logger.Errorf(ctx, "OpAMP server error: %v", err.ErrorMessage)
		},
		GetEffectiveConfig: func(ctx context.Context) (*protobufs.EffectiveConfig, error) {
			return s.getEffectiveConfig(), nil
		},
		OnMessage: s.onMessage,
	}
	callbacks.SetDefaults()

	settings := types.StartSettings{
		OpAMPServerURL: s.config.Server.Endpoint,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		InstanceUid: types.InstanceUid(s.instanceId),
		Callbacks:   callbacks,
		Capabilities: protobufs.AgentCapabilities_AgentCapabilities_AcceptsRemoteConfig |
			protobufs.AgentCapabilities_AgentCapabilities_ReportsRemoteConfig |
			protobufs.AgentCapabilities_AgentCapabilities_ReportsEffectiveConfig |
			protobufs.AgentCapabilities_AgentCapabilities_ReportsHealth |
			protobufs.AgentCapabilities_AgentCapabilities_AcceptsPackages |
			protobufs.AgentCapabilities_AgentCapabilities_ReportsPackageStatuses,
	}

	err := s.opampClient.SetAgentDescription(s.createAgentDescription())
	if err != nil {
		return fmt.Errorf("failed to set agent description: %w", err)
	}

	err = s.opampClient.SetHealth(&protobufs.ComponentHealth{Healthy: false})
	if err != nil {
		return fmt.Errorf("failed to set initial health: %w", err)
	}

	err = s.opampClient.Start(ctx, settings)
	if err != nil {
		return fmt.Errorf("failed to start OpAMP client: %w", err)
	}

	s.logger.Debugf(ctx, "OpAMP client started successfully.")
	return nil
}

// Stop gracefully stops the supervisor and managed collector.
func (s *Supervisor) Stop(ctx context.Context) error {
	s.logger.Debugf(ctx, "Supervisor shutting down...")
	s.stopCollector()
	if s.opampClient != nil {
		_ = s.opampClient.SetHealth(&protobufs.ComponentHealth{
			Healthy:   false,
			LastError: "Supervisor shutdown",
		})
		_ = s.opampClient.Stop(ctx)
	}
	return nil
}

// Upgrader returns the binary upgrader instance.
func (s *Supervisor) Upgrader() *BinaryUpgrader { return s.upgrader }

// InstanceID returns the supervisor's instance UUID.
func (s *Supervisor) InstanceID() uuid.UUID { return s.instanceId }

// AgentVersion returns the current agent version.
func (s *Supervisor) AgentVersion() string { return s.agentVersion }

// Config returns the supervisor configuration.
func (s *Supervisor) Config() *Config { return s.config }

// IsCollectorRunning returns whether the collector process is running.
func (s *Supervisor) IsCollectorRunning() bool { return s.running.Load() }

// StartedAt returns the time the collector was last started.
func (s *Supervisor) StartedAt() time.Time { return s.startedAt }

// Hostname returns the hostname of the machine.
func (s *Supervisor) Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// StartCollector starts the collector process.
func (s *Supervisor) StartCollector() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running.Load() {
		return nil
	}
	execPath := s.config.Agent.Executable
	if execPath == "" {
		return fmt.Errorf("agent executable not configured")
	}
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		return fmt.Errorf("agent executable not found: %s", execPath)
	}
	s.cmd = exec.Command(execPath, "--config", "effective.yaml")
	s.cmd.Stdout = os.Stdout
	s.cmd.Stderr = os.Stderr
	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start collector: %w", err)
	}
	s.running.Store(true)
	s.startedAt = time.Now()
	s.logger.Debugf(context.Background(), "Collector started, PID=%d", s.cmd.Process.Pid)
	go func() {
		err := s.cmd.Wait()
		s.running.Store(false)
		if err != nil {
			s.logger.Errorf(context.Background(), "Collector exited: %v", err)
		}
		_ = s.opampClient.SetHealth(&protobufs.ComponentHealth{
			Healthy:   false,
			LastError: fmt.Sprintf("Collector process exited: %v", err),
		})
	}()
	return nil
}

func (s *Supervisor) stopCollector() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cmd != nil && s.cmd.Process != nil {
		s.logger.Debugf(context.Background(), "Stopping collector, PID=%d", s.cmd.Process.Pid)
		_ = s.cmd.Process.Kill()
		_, _ = s.cmd.Process.Wait()
		s.running.Store(false)
	}
}

func (s *Supervisor) getEffectiveConfig() *protobufs.EffectiveConfig {
	return &protobufs.EffectiveConfig{
		ConfigMap: &protobufs.AgentConfigMap{
			ConfigMap: map[string]*protobufs.AgentConfigFile{
				"": {Body: []byte(composeCollectorConfig(s.instanceId.String(), s.agentVersion))},
			},
		},
	}
}

func (s *Supervisor) onMessage(ctx context.Context, msg *types.MessageData) {
	if msg.RemoteConfig != nil {
		s.logger.Debugf(ctx, "Received remote config from server.")
		_ = s.opampClient.SetRemoteConfigStatus(&protobufs.RemoteConfigStatus{
			Status: protobufs.RemoteConfigStatuses_RemoteConfigStatuses_APPLIED,
		})
	}
	if msg.PackagesAvailable != nil {
		s.logger.Debugf(ctx, "Received package offer from server.")
		err := s.upgrader.ProcessPackageOffer(ctx, msg.PackagesAvailable, s.config.Agent.Executable)
		if err != nil {
			s.logger.Errorf(ctx, "Failed to process package offer: %v", err)
		} else {
			s.logger.Debugf(ctx, "Package upgrade applied. Restarting collector...")
			s.stopCollector()
			_ = s.StartCollector()
		}
	}
	if msg.AgentIdentification != nil {
		newID, err := uuid.FromBytes(msg.AgentIdentification.NewInstanceUid)
		if err == nil {
			s.instanceId = newID
		}
	}
}

func (s *Supervisor) createAgentDescription() *protobufs.AgentDescription {
	hostname, _ := os.Hostname()
	return &protobufs.AgentDescription{
		IdentifyingAttributes: []*protobufs.KeyValue{
			{Key: "service.name", Value: &protobufs.AnyValue{Value: &protobufs.AnyValue_StringValue{StringValue: "io.opentelemetry.collector"}}},
			{Key: "service.version", Value: &protobufs.AnyValue{Value: &protobufs.AnyValue_StringValue{StringValue: s.agentVersion}}},
		},
		NonIdentifyingAttributes: []*protobufs.KeyValue{
			{Key: "host.name", Value: &protobufs.AnyValue{Value: &protobufs.AnyValue_StringValue{StringValue: hostname}}},
		},
	}
}

func composeCollectorConfig(instanceID, version string) string {
	return fmt.Sprintf(`
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318
exporters:
  debug:
    verbosity: detailed
service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [debug]
    metrics:
      receivers: [otlp]
      exporters: [debug]
    logs:
      receivers: [otlp]
      exporters: [debug]
`)
}