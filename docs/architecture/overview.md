# Architecture Overview

Build a centralized management platform for OpenTelemetry Collector agents.

Components:

1. Control Plane UI
- React
- Collector inventory
- Operations

2. Control Plane Backend
- Go
- REST APIs
- OPAMP communication
- State management

3. OPAMP Layer
- Agent communication
- Desired state management

4. Collector Fleet
- Managed OpenTelemetry collectors
