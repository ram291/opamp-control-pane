package api

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ram291/opamp-control-pane/internal/api/handlers"
	"github.com/ram291/opamp-control-pane/internal/apiuser"
	"github.com/ram291/opamp-control-pane/internal/supervisor"
)

// Server wraps the HTTP API server.
type Server struct {
	supervisor *supervisor.Supervisor
	router     chi.Router
	frontendFS fs.FS
	h          *handlers.Handlers
}

// New creates a new API server.
func New(sup *supervisor.Supervisor, frontendFS fs.FS) *Server {
	s := &Server{
		supervisor: sup,
		frontendFS: frontendFS,
		h:          handlers.New(sup),
	}
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(s.corsMiddleware)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Use(s.authMiddleware)

		// Agent endpoints (read-only: all authenticated users)
		r.With(s.requirePermission(apiuser.PermAgentList)).
			Get("/agents", s.h.ListAgents)
		r.With(s.requirePermission(apiuser.PermAgentView)).
			Get("/agents/{id}", s.h.GetAgent)

		// Upgrade endpoints
		r.With(s.requirePermission(apiuser.PermUpgradeView)).
			Get("/versions", s.h.ListVersions)

		// Config endpoints
		r.With(s.requirePermission(apiuser.PermConfigView)).
			Get("/agents/{id}/config", s.h.GetConfig)

		// Actions (config-deployer and above)
		r.With(s.requirePermission(apiuser.PermUpgradeExec)).
			Post("/agents/{id}/upgrade", s.h.UpgradeAgent)

		// User info
		r.Get("/me", s.h.CurrentUser)
	})

	// Auth endpoints
	r.Get("/auth/login", s.loginHandler)
	r.Get("/auth/callback", s.callbackHandler)

	// Serve React SPA for all non-API routes
	r.Group(func(r chi.Router) {
		r.Use(s.spaFallback)
		r.Handle("/*", s.serveFrontend())
	})

	s.router = r
	return s
}

// Handler returns the HTTP handler.
func (s *Server) Handler() http.Handler {
	return s.router
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) spaFallback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") || strings.HasPrefix(r.URL.Path, "/auth") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) serveFrontend() http.Handler {
	if s.frontendFS != nil {
		return http.FileServer(http.FS(s.frontendFS))
	}
	return http.NotFoundHandler()
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &apiuser.User{
			ID:    "dev-user",
			Email: "dev@example.com",
			Name:  "Development User",
			Roles: []apiuser.Role{apiuser.RoleAdmin},
		}
		ctx := apiuser.WithUserContext(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) requirePermission(permission apiuser.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := apiuser.GetUserFromContext(r.Context())
			if user == nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			if !user.HasPermission(permission) {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "login not configured"})
}

func (s *Server) callbackHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "callback not configured"})
}