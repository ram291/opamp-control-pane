# OpAMP Control Pane

A production-grade Open Agent Management Protocol (OpAMP) supervisor for managing OpenTelemetry Collector instances with Artifactory-based binary management, RBAC via OIDC, and a React-based management UI.

## Architecture

```
┌────────────────────────────────────────────────────────────────┐
│                    opamp-control-pane                          │
│                                                                │
│  ┌──────────────────────────────┐   ┌──────────────────────┐  │
│  │       Supervisor Core        │   │   Management API     │  │
│  │                              │   │   (REST/JSON)        │  │
│  │  • OpAMP WebSocket Client    │   │                      │  │
│  │  • Collector Process Manager │   │  • Agent management  │  │
│  │  • Binary Upgrade from       │   │  • Config management │  │
│  │    rpm                       │   │  • Upgrade triggers  │  │
│  │  • Config Composition        │   │  • Health monitoring │  │
│  │  • Health Monitoring         │   │  • RBAC (OIDC)       │  │
│  └──────────────┬───────────────┘   └──────────┬───────────┘  │
│                 │                               │              │
│                 │         Same Process          │              │
│                 └───────────────┬───────────────┘              │
│                                 │                              │
│                    ┌────────────▼────────────┐                 │
│                    │   React SPA (Embedded)  │                 │
│                    │   • Dashboard           │                 │
│                    │   • Agent Management    │                 │
│                    │   • Config Editor       │                 │
│                    │   • Upgrade Workflow    │                 │
│                    │   • RBAC-gated UI       │                 │
│                    └─────────────────────────┘                 │
└────────────────────────────────────────────────────────────────┘
```

## Features

- **OpAMP Protocol**: Full OpAMP WebSocket client for agent management
- **Compont Binary Upgrade using Supervisor RPM**: Binary upgrades using rpm to be fetched from the unix repo.
- **RBAC**: Role-based access control via OIDC (Azure AD, Okta, Keycloak)
- **Web UI**: React-based management dashboard (embedded in binary)
- **Config Management**: Compose remote, local, and dynamic config sections. Collectors configuration to be saved and pushed from the github repo.
- **Health Monitoring**: Automatic health checks and reporting
- **Single Binary**: Go binary with embedded React frontend

## Prerequisites

- Go 1.22+
- Node.js 20+
- npm

## Quick Start

```bash
# Clone the repository
git clone https://github.com/ram291/opamp-control-pane.git
cd opamp-control-pane

# Install dependencies
make deps

# Build the frontend and backend
make build

# Run
./bin/supervisor -config configs/supervisor.yaml
```

## Development

```bash
# Terminal 1: Run Go backend with hot-reload
make dev

# Terminal 2: Run React dev server with HMR
make ui-dev

# Or build everything
make all
```

## Configuration

See [configs/supervisor.example.yaml](configs/supervisor.example.yaml) for configuration options.

### Environment Variables

| Variable | Description |
|----------|-------------|
| `OIDC_CLIENT_ID` | OIDC client ID |
| `OIDC_CLIENT_SECRET` | OIDC client secret |
| `OPAMP_SERVER_URL` | OpAMP server WebSocket URL |

## RBAC Roles

| Role | Permissions |
|------|-------------|
| **admin** | Full access (agents, config, upgrades, admin, audit) |
| **config-deployer** | Manage configs and upgrades, view agents |
| **read-only** | View agents, configs, and status |

## Updating opamp-go

```bash
make update-opamp
```

This fetches the latest version from the upstream community repository.

## License

MIT
