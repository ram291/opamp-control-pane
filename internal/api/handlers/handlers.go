package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ram291/opamp-control-pane/internal/apiuser"
	"github.com/ram291/opamp-control-pane/internal/supervisor"
)

// Handlers holds the HTTP handler functions.
type Handlers struct {
	supervisor *supervisor.Supervisor
}

// New creates a new Handlers instance.
func New(sup *supervisor.Supervisor) *Handlers {
	return &Handlers{supervisor: sup}
}

// AgentDTO is the data transfer object for agent information.
type AgentDTO struct {
	ID       string `json:"id"`
	Version  string `json:"version"`
	Status   string `json:"status"`
	Uptime   string `json:"uptime"`
	Hostname string `json:"hostname"`
}

// VersionDTO represents an available collector version.
type VersionDTO struct {
	Version string `json:"version"`
	Date    string `json:"date"`
	SHA256  string `json:"sha256"`
}

// UserDTO represents the current user.
type UserDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

// ListAgents returns all managed agents.
func (h *Handlers) ListAgents(w http.ResponseWriter, r *http.Request) {
	status := "stopped"
	if h.supervisor.IsCollectorRunning() {
		status = "running"
	}

	agents := []AgentDTO{
		{
			ID:       h.supervisor.InstanceID().String(),
			Version:  h.supervisor.AgentVersion(),
			Status:   status,
			Hostname: h.supervisor.Hostname(),
		},
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"agents": agents,
	})
}

// GetAgent returns details for a specific agent.
func (h *Handlers) GetAgent(w http.ResponseWriter, r *http.Request) {
	status := "stopped"
	if h.supervisor.IsCollectorRunning() {
		status = "running"
	}

	agent := AgentDTO{
		ID:       h.supervisor.InstanceID().String(),
		Version:  h.supervisor.AgentVersion(),
		Status:   status,
		Hostname: h.supervisor.Hostname(),
	}

	json.NewEncoder(w).Encode(agent)
}

// ListVersions returns available collector versions from Artifactory.
func (h *Handlers) ListVersions(w http.ResponseWriter, r *http.Request) {
	versions, err := h.supervisor.Upgrader().AvailableVersions(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to list versions"}`, http.StatusInternalServerError)
		return
	}

	var dtos []VersionDTO
	for _, v := range versions {
		dtos = append(dtos, VersionDTO{
			Version: v.Version,
			Date:    v.ReleaseDate,
			SHA256:  v.SHA256,
		})
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"versions": dtos,
	})
}

// GetConfig returns the effective config of the agent.
func (h *Handlers) GetConfig(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"config": map[string]interface{}{
			"content": "effective config placeholder",
			"source":  "local",
		},
	})
}

// UpgradeAgent triggers a collector binary upgrade.
func (h *Handlers) UpgradeAgent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "upgrade_triggered",
		"version": req.Version,
	})
}

// CurrentUser returns the authenticated user's information.
func (h *Handlers) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := apiuser.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var permissions []string
	for _, role := range user.Roles {
		if perms, ok := apiuser.RolePermissions[role]; ok {
			for _, p := range perms {
				permissions = append(permissions, string(p))
			}
		}
	}

	roleStrs := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roleStrs[i] = string(r)
	}

	json.NewEncoder(w).Encode(UserDTO{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Roles:       roleStrs,
		Permissions: permissions,
	})
}