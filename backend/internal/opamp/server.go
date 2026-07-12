package opamp

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/ram291/opamp-control-pane/backend/internal/agent"
	"github.com/ram291/opamp-control-pane/backend/internal/storage"
)

// Server handles the backend side of the OPAMP registration workflow.
type Server struct {
	registry *agent.Registry
	repo     *storage.Repository
}

func NewServer() *Server {
	log.Println("Initializing OPAMP server")

	return &Server{
		registry: agent.NewRegistry(),
		repo:     storage.NewRepository(),
	}
}

func (s *Server) Start() {
	log.Println("OPAMP server started")
}

// HandleRegistration implements the ADR-004 registration handshake.
func (s *Server) HandleRegistration(_ context.Context, msg *protobufs.AgentToServer) *protobufs.ServerToAgent {
	if msg == nil {
		return &protobufs.ServerToAgent{}
	}

	instanceUID := string(msg.InstanceUid)
	if instanceUID == "" || (msg.Flags&uint64(protobufs.AgentToServerFlags_AgentToServerFlags_RequestInstanceUid)) != 0 {
		instanceUID = s.generateInstanceUID()
	}

	registered := &agent.Agent{
		ID:               instanceUID,
		InstanceUID:      instanceUID,
		Capabilities:     msg.Capabilities,
		IdentifyingAttrs: extractAttributes(msg.AgentDescription),
	}

	s.registry.RegisterAgent(registered)
	s.repo.SaveAgent(registered.ID, registered.InstanceUID)

	response := &protobufs.ServerToAgent{
		InstanceUid: []byte(instanceUID),
		Capabilities: uint64(protobufs.ServerCapabilities_ServerCapabilities_AcceptsStatus | protobufs.ServerCapabilities_ServerCapabilities_OffersRemoteConfig),
	}

	if (msg.Flags & uint64(protobufs.AgentToServerFlags_AgentToServerFlags_RequestInstanceUid)) != 0 || len(msg.InstanceUid) == 0 {
		response.AgentIdentification = &protobufs.AgentIdentification{NewInstanceUid: []byte(instanceUID)}
	}

	return response
}

func (s *Server) generateInstanceUID() string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

func extractAttributes(description *protobufs.AgentDescription) map[string]string {
	if description == nil {
		return map[string]string{}
	}

	attrs := make(map[string]string, len(description.IdentifyingAttributes))
	for _, attr := range description.IdentifyingAttributes {
		if attr != nil && attr.Key != "" {
			attrs[attr.Key] = attr.Value
		}
	}
	return attrs
}

