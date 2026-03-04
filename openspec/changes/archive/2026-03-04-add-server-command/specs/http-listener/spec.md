## MODIFIED Requirements

### Requirement: Server responds with hello world
The server SHALL use the host-routing handler as its HTTP handler instead of responding with a static "hello world" body. All incoming requests SHALL be dispatched through the host-routing logic.

#### Scenario: Any HTTP request with matching host
- **WHEN** a client sends an HTTP request with a matching `Host` header
- **THEN** the server SHALL forward the request to the configured upstream via the host-routing handler

#### Scenario: Any HTTP request with non-matching host
- **WHEN** a client sends an HTTP request with a non-matching `Host` header
- **THEN** the server SHALL respond with HTTP 404

### Requirement: Server binds to all interfaces
The server SHALL bind to `127.0.0.1` instead of `0.0.0.0`, since the proxy only serves localhost traffic for `.localhost` domains.

#### Scenario: Default bind address
- **WHEN** the server starts with no CLI arguments
- **THEN** it SHALL bind to `127.0.0.1:7890`
