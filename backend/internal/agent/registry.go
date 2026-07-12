package agent

import "sync"

// Agent represents a registered OPAMP agent in the control plane.
type Agent struct {
	ID               string
	InstanceUID      string
	Capabilities     uint64
	IdentifyingAttrs map[string]string
}

// Registry stores agents known to the control plane.
type Registry struct {
	mu      sync.RWMutex
	agents  map[string]*Agent
}

func NewRegistry() *Registry {
	return &Registry{agents: make(map[string]*Agent)}
}

// RegisterAgent stores or updates an agent registration.
func (r *Registry) RegisterAgent(agent *Agent) {
	if agent == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if agent.ID == "" {
		agent.ID = agent.InstanceUID
	}

	r.agents[agent.ID] = agent
}

// ListAgents returns all currently known agents.
func (r *Registry) ListAgents() []*Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agents := make([]*Agent, 0, len(r.agents))
	for _, agent := range r.agents {
		agents = append(agents, agent)
	}
	return agents
}

