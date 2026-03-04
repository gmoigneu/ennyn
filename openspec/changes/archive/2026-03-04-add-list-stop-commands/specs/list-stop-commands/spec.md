## ADDED Requirements

### Requirement: List running instances
The `ennyn list` command SHALL display all active Ennyn server instances. For each instance, it SHALL show: hostname, PID, proxy port, app port, app command, and how long it has been running.

#### Scenario: Multiple instances running
- **WHEN** `ennyn list` is run and two instances are active (myapp and api)
- **THEN** the output SHALL list both instances with their metadata

#### Scenario: No instances running
- **WHEN** `ennyn list` is run and no instances are active
- **THEN** the output SHALL print a message indicating no running instances

#### Scenario: Stale entries cleaned up
- **WHEN** `ennyn list` is run and some state files reference dead processes
- **THEN** the stale files SHALL be removed and only live instances SHALL be displayed

### Requirement: Stop a running instance by hostname
The `ennyn stop <host>` command SHALL send SIGTERM to the Ennyn server instance managing the specified hostname.

#### Scenario: Stop a running instance
- **WHEN** `ennyn stop myapp` is run and myapp is active
- **THEN** SIGTERM SHALL be sent to the instance's PID and a confirmation message SHALL be printed

#### Scenario: Stop a non-existent instance
- **WHEN** `ennyn stop myapp` is run and no instance for myapp exists
- **THEN** the command SHALL print an error message and exit with a non-zero code

#### Scenario: Stop with stale state file
- **WHEN** `ennyn stop myapp` is run and the state file exists but the PID is not running
- **THEN** the stale file SHALL be removed and an appropriate message SHALL be printed

### Requirement: Stop requires a hostname argument
The `ennyn stop` command SHALL require exactly one hostname argument. If none is provided, it SHALL print a usage message and exit with a non-zero code.

#### Scenario: No hostname provided
- **WHEN** `ennyn stop` is run with no arguments
- **THEN** the command SHALL print usage and exit with a non-zero code
