## ADDED Requirements

### Requirement: State file creation on server start
When an `ennyn server` instance starts successfully, it SHALL write a JSON state file at `<run-dir>/<host>.json`. The state file SHALL contain: `pid` (int), `host` (string), `proxyPort` (int), `appPort` (int), `appCmd` (string), `startedAt` (RFC 3339 timestamp).

#### Scenario: Server writes state file on start
- **WHEN** `ennyn server myapp 3000 npm run dev` starts successfully
- **THEN** a file `<run-dir>/myapp.json` SHALL be created with the instance metadata

### Requirement: State file removal on server shutdown
When an `ennyn server` instance shuts down (via signal or child exit), it SHALL remove its state file.

#### Scenario: Server removes state file on SIGTERM
- **WHEN** a running instance receives SIGTERM and shuts down
- **THEN** its state file SHALL be deleted from the run directory

#### Scenario: Server removes state file on child exit
- **WHEN** the child process exits and the server shuts down
- **THEN** its state file SHALL be deleted from the run directory

### Requirement: State directory location
The state directory SHALL be `$XDG_CONFIG_HOME/ennyn/run/` if `$XDG_CONFIG_HOME` is set, otherwise `~/.config/ennyn/run/`. The directory SHALL be created automatically if it does not exist.

#### Scenario: Default state directory
- **WHEN** `$XDG_CONFIG_HOME` is not set
- **THEN** state files SHALL be stored at `~/.config/ennyn/run/`

#### Scenario: Custom XDG path
- **WHEN** `$XDG_CONFIG_HOME` is set to `/custom/config`
- **THEN** state files SHALL be stored at `/custom/config/ennyn/run/`

### Requirement: Stale state file detection
When reading state files, the system SHALL verify that the PID in the file corresponds to a running process. If the process is not running, the state file SHALL be considered stale.

#### Scenario: Process is running
- **WHEN** a state file contains PID 1234 and process 1234 is alive
- **THEN** the instance SHALL be reported as active

#### Scenario: Process is not running
- **WHEN** a state file contains PID 1234 and process 1234 is not running
- **THEN** the state file SHALL be considered stale and cleaned up automatically

### Requirement: Duplicate host prevention
If a state file already exists for a hostname and the PID is still alive, `ennyn server` SHALL refuse to start and exit with an error message indicating the host is already in use.

#### Scenario: Host already running
- **WHEN** `ennyn server myapp 3000 npm run dev` is started and `myapp` is already registered with a live PID
- **THEN** the server SHALL print an error and exit with a non-zero code

#### Scenario: Host has stale state file
- **WHEN** `ennyn server myapp 3000 npm run dev` is started and `myapp` has a state file with a dead PID
- **THEN** the stale file SHALL be removed and the server SHALL start normally
