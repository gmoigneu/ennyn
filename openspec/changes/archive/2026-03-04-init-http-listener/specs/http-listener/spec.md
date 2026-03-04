## ADDED Requirements

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
The server SHALL respond to all HTTP requests with a plain-text body containing "hello world".

#### Scenario: Any HTTP request
- **WHEN** a client sends an HTTP request to any path
- **THEN** the server SHALL respond with status 200, Content-Type `text/plain`, and body `hello world`

### Requirement: Server binds to all interfaces
The server SHALL bind to `0.0.0.0` by default, making it accessible from any network interface.

#### Scenario: Default bind address
- **WHEN** the server starts with no CLI arguments
- **THEN** it SHALL be reachable on `127.0.0.1:7890` and on the host's LAN IP at port 7890

### Requirement: Graceful shutdown
The server SHALL shut down gracefully when receiving SIGINT or SIGTERM, draining in-flight requests before exiting.

#### Scenario: SIGINT received
- **WHEN** the server receives SIGINT while handling requests
- **THEN** it SHALL stop accepting new connections, finish in-flight requests (up to a 5-second timeout), and exit with code 0

#### Scenario: SIGTERM received
- **WHEN** the server receives SIGTERM
- **THEN** it SHALL behave identically to receiving SIGINT
