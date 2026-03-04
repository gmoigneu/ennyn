### Requirement: Server binds to configurable port
The server SHALL listen for HTTP connections on port 7890 by default. The port MUST be configurable via `-p` or `-port` CLI flags.

#### Scenario: Default port
- **WHEN** the server starts with no CLI arguments
- **THEN** it SHALL bind to `0.0.0.0:7890`

#### Scenario: Custom port via -p flag
- **WHEN** the server starts with `-p 3000`
- **THEN** it SHALL bind to `0.0.0.0:3000`

#### Scenario: Custom port via -port flag
- **WHEN** the server starts with `-port 3000`
- **THEN** it SHALL bind to `0.0.0.0:3000`

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

### Requirement: Graceful shutdown
The server SHALL shut down gracefully when receiving SIGINT or SIGTERM, draining in-flight requests before exiting. The server SHALL remove its state file from the registry during shutdown.

#### Scenario: SIGINT received
- **WHEN** the server receives SIGINT while handling requests
- **THEN** it SHALL stop accepting new connections, finish in-flight requests (up to a 5-second timeout), remove its state file, and exit with code 0

#### Scenario: SIGTERM received
- **WHEN** the server receives SIGTERM
- **THEN** it SHALL behave identically to receiving SIGINT
