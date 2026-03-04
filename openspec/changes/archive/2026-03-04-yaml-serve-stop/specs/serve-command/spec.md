## ADDED Requirements

### Requirement: Serve subcommand starts all services from config
The `ennyn serve` command SHALL read `ennyn.yml` from the current directory and start an `ennyn server` background process for each service defined in the config.

#### Scenario: Start services from config
- **WHEN** the user runs `ennyn serve` and `ennyn.yml` defines two services (myapp on port 3000, api on port 8080)
- **THEN** two background `ennyn server` processes SHALL be started, one for each service

#### Scenario: Config file not found
- **WHEN** the user runs `ennyn serve` and no `ennyn.yml` exists in the current directory
- **THEN** the command SHALL print an error message and exit with a non-zero code

### Requirement: Serve reports startup results
After starting all services, `ennyn serve` SHALL print a summary showing each service's hostname, port, and whether it started successfully.

#### Scenario: All services start successfully
- **WHEN** `ennyn serve` starts two services and both succeed
- **THEN** the output SHALL list both services as started with their hostnames and ports

#### Scenario: One service fails to start
- **WHEN** `ennyn serve` starts two services and one fails (e.g., invalid command)
- **THEN** the output SHALL show the successful service as started and the failed service with an error, and the exit code SHALL be non-zero

### Requirement: Serve skips already-running services
If a service defined in the config is already running (detected via the instance registry), `ennyn serve` SHALL skip it and print a message indicating it is already running.

#### Scenario: Service already running
- **WHEN** `ennyn serve` is run and `myapp` is already registered with a live PID
- **THEN** the command SHALL skip starting `myapp` and print a message that it is already running

### Requirement: Serve accepts proxy port flag
The `ennyn serve` command SHALL accept `-p` or `-port` flags to set the proxy listening port for all started services. The default SHALL be 7890.

#### Scenario: Custom proxy port
- **WHEN** the user runs `ennyn serve -p 9000`
- **THEN** all started server instances SHALL use proxy port 9000

### Requirement: Serve starts server processes in background
Each `ennyn server` process started by `ennyn serve` SHALL be a detached background process that continues running after `ennyn serve` exits. The server processes SHALL NOT inherit `ennyn serve`'s stdin.

#### Scenario: Serve exits after starting
- **WHEN** `ennyn serve` finishes starting all services
- **THEN** it SHALL exit and the server processes SHALL continue running independently
