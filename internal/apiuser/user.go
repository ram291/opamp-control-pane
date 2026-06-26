package apiuser

import "context"

// Role defines a named set of permissions.
type Role string

const (
	RoleAdmin          Role = "admin"
	RoleConfigDeployer Role = "config-deployer"
	RoleReadOnly       Role = "read-only"
)

// Permission defines a specific action that can be performed.
type Permission string

const (
	PermAgentList     Permission = "agent:list"
	PermAgentView     Permission = "agent:view"
	PermUpgradeView   Permission = "upgrade:view"
	PermUpgradeExec   Permission = "upgrade:execute"
	PermConfigView    Permission = "config:view"
	PermConfigEdit    Permission = "config:edit"
	PermConfigDeploy  Permission = "config:deploy"
	PermAdminUsers    Permission = "admin:users"
	PermAdminSettings Permission = "admin:settings"
	PermAuditLog      Permission = "admin:audit"
)

// RolePermissions maps roles to their allowed permissions.
var RolePermissions = map[Role][]Permission{
	RoleAdmin: {
		PermAgentList, PermAgentView,
		PermUpgradeView, PermUpgradeExec,
		PermConfigView, PermConfigEdit, PermConfigDeploy,
		PermAdminUsers, PermAdminSettings, PermAuditLog,
	},
	RoleConfigDeployer: {
		PermAgentList, PermAgentView,
		PermUpgradeView, PermUpgradeExec,
		PermConfigView, PermConfigEdit, PermConfigDeploy,
	},
	RoleReadOnly: {
		PermAgentList, PermAgentView,
		PermUpgradeView, PermConfigView,
	},
}

// User represents an authenticated user with roles.
type User struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Name  string   `json:"name"`
	Roles []Role   `json:"roles"`
}

// HasRole checks if the user has a specific role.
func (u *User) HasRole(role Role) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasPermission checks if the user has a specific permission.
func (u *User) HasPermission(permission Permission) bool {
	for _, role := range u.Roles {
		if perms, ok := RolePermissions[role]; ok {
			for _, p := range perms {
				if p == permission {
					return true
				}
			}
		}
	}
	return false
}

type contextKey string

const userContextKey contextKey = "user"

// WithUserContext injects a user into the context.
func WithUserContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetUserFromContext extracts the user from the context.
func GetUserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		return nil
	}
	return user
}