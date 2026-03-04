## MODIFIED Requirements

### Requirement: Graceful shutdown
The server SHALL shut down gracefully when receiving SIGINT or SIGTERM, draining in-flight requests before exiting. The server SHALL remove its state file from the registry during shutdown.

#### Scenario: SIGINT received
- **WHEN** the server receives SIGINT while handling requests
- **THEN** it SHALL stop accepting new connections, finish in-flight requests (up to a 5-second timeout), remove its state file, and exit with code 0

#### Scenario: SIGTERM received
- **WHEN** the server receives SIGTERM
- **THEN** it SHALL behave identically to receiving SIGINT
