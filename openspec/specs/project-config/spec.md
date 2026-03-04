### Requirement: Config file format
The `ennyn.yml` file SHALL contain a top-level `services` key whose value is a list of service objects. Each service object SHALL have the keys `host` (string), `port` (integer), and `command` (string).

#### Scenario: Valid config with two services
- **WHEN** `ennyn.yml` contains:
  ```yaml
  services:
    - host: myapp
      port: 3000
      command: npm run dev
    - host: api
      port: 8080
      command: go run ./cmd/api
  ```
- **THEN** the parser SHALL return two service entries with the correct host, port, and command values

#### Scenario: Empty services list
- **WHEN** `ennyn.yml` contains `services: []`
- **THEN** the parser SHALL return an empty list and no error

### Requirement: Config file validation
The parser SHALL validate each service entry. The `host` field SHALL follow the same rules as the server command host argument (lowercase alphanumeric and hyphens, cannot start or end with hyphen). The `port` field SHALL be an integer in range 1–65535. The `command` field SHALL be a non-empty string.

#### Scenario: Invalid host in config
- **WHEN** a service entry has `host: "MyApp"`
- **THEN** the parser SHALL return an error indicating the host is invalid

#### Scenario: Port out of range
- **WHEN** a service entry has `port: 99999`
- **THEN** the parser SHALL return an error indicating the port is out of range

#### Scenario: Missing command
- **WHEN** a service entry has an empty `command` field
- **THEN** the parser SHALL return an error indicating the command is required

### Requirement: Duplicate host detection
The parser SHALL reject config files containing duplicate host values.

#### Scenario: Duplicate hosts
- **WHEN** `ennyn.yml` contains two services with `host: myapp`
- **THEN** the parser SHALL return an error indicating duplicate host "myapp"

### Requirement: Config file location
The config file SHALL be read from `ennyn.yml` in the current working directory. No upward directory traversal SHALL be performed.

#### Scenario: Config file exists
- **WHEN** `ennyn.yml` exists in the current directory
- **THEN** the parser SHALL read and parse it

#### Scenario: Config file not found
- **WHEN** `ennyn.yml` does not exist in the current directory
- **THEN** the system SHALL return an error indicating the config file was not found

### Requirement: Command string parsing
The `command` field SHALL be parsed into an executable and arguments by splitting on whitespace (shell-style splitting is not required — no quote handling).

#### Scenario: Command with arguments
- **WHEN** a service entry has `command: "python -m flask run"`
- **THEN** the parser SHALL split it into executable `python` and arguments `["-m", "flask", "run"]`
