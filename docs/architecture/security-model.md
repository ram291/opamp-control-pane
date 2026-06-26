# Security Model

Authentication:

OIDC

Flow:

User -> React -> Identity Provider -> JWT -> Backend


Authorization:

RBAC controls actions.

Viewer:
- View collectors

Operator:
- Update configuration

Administrator:
- Upgrade binaries


Audit operations:
- User
- Action
- Target
- Timestamp
