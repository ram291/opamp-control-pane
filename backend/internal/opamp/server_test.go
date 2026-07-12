package opamp
package opamp

import (
	"context"
	"testing"

	"github.com/open-telemetry/opamp-go/protobufs"
)

func TestHandleRegistrationRegistersAgentAndResponds(t *testing.T) {
	srv := NewServer()

	msg := &protobufs.AgentToServer{
		InstanceUid: []byte("agent-1"),
		AgentDescription: &protobufs.AgentDescription{},
		Capabilities: uint64(protobufs.AgentCapabilities_AgentCapabilities_AcceptsRemoteConfig),
	}

	response := srv.HandleRegistration(context.Background(), msg)
	if response == nil {
		t.Fatal("expected a registration response")
	}

	if string(response.InstanceUid) != string(msg.InstanceUid) {
		t.Fatalf("expected instance uid %q, got %q", msg.InstanceUid, response.InstanceUid)
	}

	if response.Capabilities == 0 {
		t.Fatal("expected server capabilities to be advertised")
	}

	if len(srv.registry.ListAgents()) != 1 {
		t.Fatalf("expected one registered agent, got %d", len(srv.registry.ListAgents()))
	}
}

func TestHandleRegistrationAssignsNewInstanceIDWhenRequested(t *testing.T) {
	srv := NewServer()

	msg := &protobufs.AgentToServer{
		Flags: uint64(protobufs.AgentToServerFlags_RequestInstanceUid),
		AgentDescription: &protobufs.AgentDescription{},
	}

	response := srv.HandleRegistration(context.Background(), msg)
	if response == nil {
		t.Fatal("expected a registration response")
	}

	if len(response.InstanceUid) == 0 {
		t.Fatal("expected server to assign an instance uid")
	}

	if response.AgentIdentification == nil || len(response.AgentIdentification.NewInstanceUid) == 0 {
		t.Fatal("expected server to provide agent identification")
	}
}
