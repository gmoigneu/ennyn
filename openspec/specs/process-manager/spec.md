## ADDED Requirements

### Requirement: Start app command as a child process
The proxy SHALL start the app command as a child process using direct exec (no shell interpretation). The child process SHALL inherit the proxy's stdout and stderr.

#### Scenario: App starts successfully
- **WHEN** `ennyn server myapp 3000 node server.js` is run and `node` is in PATH
- **THEN** the proxy SHALL start `node server.js` as a child process and the app's output SHALL appear in the proxy's terminal

#### Scenario: App command not found
- **WHEN** the app command does not exist in PATH
- **THEN** the proxy SHALL exit with code 1 and print an error indicating the command was not found

### Requirement: Stop child process on proxy shutdown
When the proxy receives SIGINT or SIGTERM, it SHALL send SIGTERM to the child process and wait up to 5 seconds for it to exit. If the child does not exit within 5 seconds, the proxy SHALL send SIGKILL.

#### Scenario: Graceful child shutdown
- **WHEN** the proxy receives SIGINT and the child exits within 5 seconds of receiving SIGTERM
- **THEN** the proxy SHALL wait for the child to exit and then shut down itself

#### Scenario: Child does not respond to SIGTERM
- **WHEN** the proxy receives SIGINT and the child does not exit within 5 seconds of receiving SIGTERM
- **THEN** the proxy SHALL send SIGKILL to the child process and then shut down

### Requirement: Proxy exits when child process exits unexpectedly
If the child process exits on its own (before the proxy receives a shutdown signal), the proxy SHALL log the child's exit status and shut down gracefully.

#### Scenario: Child crashes
- **WHEN** the child process exits with a non-zero exit code
- **THEN** the proxy SHALL log a message including the exit code and shut down with the same exit code

#### Scenario: Child exits cleanly
- **WHEN** the child process exits with code 0
- **THEN** the proxy SHALL log a message and shut down with code 0

### Requirement: Child process runs in its own process group
The child process SHALL be started in its own process group so that SIGTERM can be sent to the entire group, covering any sub-processes the app may have spawned.

#### Scenario: App spawns sub-processes
- **WHEN** the app command spawns additional child processes and the proxy shuts down
- **THEN** SIGTERM SHALL be sent to the entire process group of the app
