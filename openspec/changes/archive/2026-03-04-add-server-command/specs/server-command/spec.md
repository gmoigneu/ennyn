## ADDED Requirements

### Requirement: CLI accepts server subcommand with host, port, and command arguments
The CLI SHALL accept the invocation `ennyn server <host> <app-port> <app-cmd> [app-args...]`. The `<host>` argument defines the hostname prefix, `<app-port>` is the port the app listens on, and `<app-cmd> [app-args...]` is the command to start the application.

#### Scenario: Valid invocation
- **WHEN** the user runs `ennyn server myapp 3000 node server.js`
- **THEN** ennyn SHALL parse `myapp` as the host, `3000` as the app port, and `node server.js` as the app command

#### Scenario: Invocation with multiple app arguments
- **WHEN** the user runs `ennyn server api 8080 python -m flask run`
- **THEN** ennyn SHALL parse `api` as the host, `8080` as the app port, and `python -m flask run` as the app command with arguments

### Requirement: CLI validates arguments
The CLI SHALL validate that exactly 3 or more positional arguments follow `server`. It SHALL validate that `<app-port>` is a valid port number (1–65535).

#### Scenario: Missing arguments
- **WHEN** the user runs `ennyn server myapp`
- **THEN** ennyn SHALL exit with code 1 and print a usage message to stderr

#### Scenario: Invalid port
- **WHEN** the user runs `ennyn server myapp notanumber node server.js`
- **THEN** ennyn SHALL exit with code 1 and print an error indicating the port is invalid

#### Scenario: Port out of range
- **WHEN** the user runs `ennyn server myapp 99999 node server.js`
- **THEN** ennyn SHALL exit with code 1 and print an error indicating the port is out of range

### Requirement: CLI supports proxy port flag
The `server` subcommand SHALL accept `-p` or `-port` flags to set the proxy's listening port, defaulting to 7890.

#### Scenario: Default proxy port
- **WHEN** the user runs `ennyn server myapp 3000 node server.js`
- **THEN** the proxy SHALL listen on port 7890

#### Scenario: Custom proxy port
- **WHEN** the user runs `ennyn server -p 8000 myapp 3000 node server.js`
- **THEN** the proxy SHALL listen on port 8000

### Requirement: Host argument accepts alphanumeric and hyphen characters
The `<host>` argument SHALL only contain lowercase alphanumeric characters and hyphens. It SHALL NOT start or end with a hyphen.

#### Scenario: Valid host
- **WHEN** the user runs `ennyn server my-app 3000 node server.js`
- **THEN** ennyn SHALL accept `my-app` as a valid host

#### Scenario: Invalid host with uppercase
- **WHEN** the user runs `ennyn server MyApp 3000 node server.js`
- **THEN** ennyn SHALL exit with code 1 and print an error indicating the host is invalid

#### Scenario: Invalid host starting with hyphen
- **WHEN** the user runs `ennyn server -app 3000 node server.js`
- **THEN** ennyn SHALL exit with code 1 and print an error indicating the host is invalid
