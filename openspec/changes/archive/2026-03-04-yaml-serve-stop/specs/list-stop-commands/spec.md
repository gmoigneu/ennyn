## MODIFIED Requirements

### Requirement: Stop requires a hostname argument
The `ennyn stop` command SHALL accept either a hostname argument or no arguments. When invoked with no arguments AND an `ennyn.yml` file exists in the current directory, it SHALL stop all services defined in the config. When invoked with no arguments AND no `ennyn.yml` exists, it SHALL print a usage message and exit with a non-zero code.

#### Scenario: No hostname provided with config file
- **WHEN** `ennyn stop` is run with no arguments and `ennyn.yml` exists with services myapp and api
- **THEN** the command SHALL stop all services defined in the config

#### Scenario: No hostname provided without config file
- **WHEN** `ennyn stop` is run with no arguments and no `ennyn.yml` exists
- **THEN** the command SHALL print a usage message and exit with a non-zero code

#### Scenario: Hostname provided
- **WHEN** `ennyn stop myapp` is run
- **THEN** the command SHALL stop only the specified instance (existing behavior)

#### Scenario: Config service not running
- **WHEN** `ennyn stop` is run with no arguments and `ennyn.yml` defines service myapp but myapp is not running
- **THEN** the command SHALL print a message that myapp is not running and continue stopping other services
