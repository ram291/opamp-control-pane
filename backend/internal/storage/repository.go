package storage

import "sync"

// Repository persists control-plane state for agents.
type Repository struct {
	mu      sync.RWMutex
	agents  map[string]string
}

func NewRepository() *Repository {
	return &Repository{agents: make(map[string]string)}
}

// SaveAgent saves the current agent identifier state.
func (r *Repository) SaveAgent(id, instanceUID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.agents == nil {
		r.agents = make(map[string]string)
	}
	r.agents[id] = instanceUID
}

// GetAgent returns the stored instance UID for an agent id.
func (r *Repository) GetAgent(id string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.agents[id]
	return value, ok
}

