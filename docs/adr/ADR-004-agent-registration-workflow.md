# ADR-004 Agent Registration Workflow

Decision:

Use the standard OPAMP registration workflow for onboarding agents into the control plane.

Reason:

The OPAMP specification already defines a native registration and lifecycle handshake for agents.
This approach provides:

- A standard first-contact flow for agent discovery
- Stable agent identity through instance_uid
- Immediate capability negotiation
- Consistent behavior across reconnects and duplicate detection
- No need for a separate custom registration protocol

Design:

1. Agent connects to the OPAMP server.
2. The agent sends an initial AgentToServer status report immediately after connection.
3. The status report includes:
   - instance_uid
   - agent_description with identifying attributes
   - capabilities supported by the agent
   - initial health and status information
4. The server creates or updates the agent record in the control plane using the reported identity.
5. The server responds with a ServerToAgent message that acknowledges receipt and may include:
   - server capabilities
   - agent_identification if the server needs to assign or override the instance_uid
   - initial configuration or connection settings if applicable
6. The first successful status report plus server acknowledgment is treated as completed registration.
7. Subsequent reconnects reuse the same instance_uid when possible.
8. If a duplicate or conflicting instance_uid is detected, the server may issue a new one and the agent must adopt it for future communication.

Consequences:

- Agent registration becomes protocol-native and interoperable with standard OPAMP implementations.
- The control plane can reliably track agents across restarts and reconnects.
- Server-side persistence of agent identity and status becomes a required part of the design.
- The design depends on the agent reporting accurate identifying attributes and capabilities.

Alternatives considered:

- A custom REST-based registration endpoint.
  - Rejected because it adds a parallel onboarding flow and diverges from the standard protocol.
- Using hostname or IP address as the primary identifier.
  - Rejected because it is less stable across restarts, migrations, and container reuse.
