### Requirement: CA generation
The `ennyn trust` command SHALL generate an ECDSA P-256 root CA certificate and private key if none exist. The CA certificate SHALL be valid for 10 years. The CA certificate SHALL have `IsCA: true` and `BasicConstraintsValid: true`. The Subject Common Name SHALL be "Ennyn Local CA".

#### Scenario: First run generates CA
- **WHEN** `ennyn trust` is run and no CA exists at the storage path
- **THEN** the system SHALL generate a new CA key pair and certificate, write them to the storage directory, and proceed to trust store installation

#### Scenario: Subsequent run reuses existing CA
- **WHEN** `ennyn trust` is run and a valid CA already exists at the storage path
- **THEN** the system SHALL reuse the existing CA and proceed to trust store installation without generating a new one

### Requirement: CA storage location
The CA key and certificate SHALL be stored at `$XDG_CONFIG_HOME/ennyn/ca/` if `$XDG_CONFIG_HOME` is set, otherwise at `~/.config/ennyn/ca/`. The directory SHALL be created with 0700 permissions. The key file (`ca.key`) SHALL be written with 0600 permissions. The certificate file SHALL be named `ca.crt`.

#### Scenario: Default storage path
- **WHEN** `$XDG_CONFIG_HOME` is not set
- **THEN** the CA files SHALL be stored at `~/.config/ennyn/ca/ca.crt` and `~/.config/ennyn/ca/ca.key`

#### Scenario: Custom XDG path
- **WHEN** `$XDG_CONFIG_HOME` is set to `/custom/config`
- **THEN** the CA files SHALL be stored at `/custom/config/ennyn/ca/ca.crt` and `/custom/config/ennyn/ca/ca.key`

### Requirement: Wildcard leaf certificate for .localhost
The `ennyn trust` command SHALL generate a wildcard leaf certificate signed by the local CA. The certificate SHALL include the following Subject Alternative Names: `*.localhost`, `localhost`, `127.0.0.1`, and `::1`. The leaf certificate SHALL use ECDSA P-256 and be valid for 825 days. The certificate SHALL be stored as `localhost.crt` and the key as `localhost.key` in the same directory as the CA files.

#### Scenario: Leaf cert generated alongside CA
- **WHEN** `ennyn trust` is run and no leaf certificate exists
- **THEN** the system SHALL generate a wildcard leaf certificate signed by the CA, with SANs `*.localhost`, `localhost`, `127.0.0.1`, and `::1`, and store it at the CA storage path

#### Scenario: Leaf cert regenerated when CA is regenerated
- **WHEN** `ennyn trust` is run and the CA was just regenerated (e.g., after uninstall + reinstall)
- **THEN** the system SHALL also regenerate the leaf certificate signed by the new CA

#### Scenario: Existing leaf cert reused
- **WHEN** `ennyn trust` is run and both the CA and leaf certificate already exist and the leaf cert was signed by the current CA
- **THEN** the system SHALL reuse the existing leaf certificate without regenerating

#### Scenario: HTTPS works for any .localhost subdomain
- **WHEN** the proxy serves HTTPS using the leaf certificate and a client connects to `myapp.localhost`
- **THEN** the TLS handshake SHALL succeed without certificate warnings (assuming the CA is trusted)

### Requirement: Trust store installation
The `ennyn trust` command SHALL install the CA certificate into the operating system's trust store. On macOS, it SHALL use the System Keychain. On Linux, it SHALL use the system ca-certificates mechanism. On Windows, it SHALL use the system certificate store.

#### Scenario: macOS trust installation
- **WHEN** `ennyn trust` is run on macOS
- **THEN** the CA certificate SHALL be added to the System Keychain and trusted for SSL

#### Scenario: Linux trust installation
- **WHEN** `ennyn trust` is run on Linux
- **THEN** the CA certificate SHALL be installed via the system ca-certificates directory and `update-ca-certificates` (or equivalent for the distribution)

#### Scenario: Windows trust installation
- **WHEN** `ennyn trust` is run on Windows
- **THEN** the CA certificate SHALL be added to the Windows system certificate store

### Requirement: WSL detection and cross-environment trust
When running under WSL, `ennyn trust` SHALL detect the WSL environment and offer to install the CA certificate in the Windows host trust store in addition to the Linux trust store.

#### Scenario: WSL detected, user accepts Windows install
- **WHEN** `ennyn trust` is run under WSL and the user accepts the Windows trust prompt
- **THEN** the CA certificate SHALL be installed in both the Linux trust store and the Windows host trust store

#### Scenario: WSL detected, user declines Windows install
- **WHEN** `ennyn trust` is run under WSL and the user declines the Windows trust prompt
- **THEN** the CA certificate SHALL be installed only in the Linux trust store

#### Scenario: Not WSL
- **WHEN** `ennyn trust` is run outside WSL
- **THEN** no WSL-specific prompt SHALL be shown

### Requirement: Trust uninstallation
The `ennyn trust --uninstall` command SHALL remove the CA certificate from the OS trust store and delete the CA files from the storage directory.

#### Scenario: Uninstall removes trust and files
- **WHEN** `ennyn trust --uninstall` is run and a CA exists
- **THEN** the CA certificate SHALL be removed from the OS trust store and the CA files SHALL be deleted from disk

#### Scenario: Uninstall when no CA exists
- **WHEN** `ennyn trust --uninstall` is run and no CA exists
- **THEN** the command SHALL print a message indicating no CA was found and exit successfully

### Requirement: Post-install user instructions
After successful trust installation, the command SHALL print a summary of what was done and any next steps the user needs to take.

#### Scenario: Successful install output
- **WHEN** `ennyn trust` completes successfully
- **THEN** the output SHALL include: the CA storage path, which trust stores were updated, and a note that HTTPS is now ready for use with Ennyn

### Requirement: Subcommand dispatch
The binary SHALL support subcommands. `ennyn trust` SHALL invoke the CA trust workflow. Running `ennyn` with no subcommand (or with flags like `-p`) SHALL start the HTTP server as before.

#### Scenario: No subcommand starts server
- **WHEN** `ennyn` is run with no arguments or with `-p <port>`
- **THEN** the HTTP server SHALL start as before

#### Scenario: Trust subcommand
- **WHEN** `ennyn trust` is run
- **THEN** the CA trust workflow SHALL execute

#### Scenario: Unknown subcommand
- **WHEN** `ennyn foobar` is run
- **THEN** the system SHALL print a usage message listing available commands and exit with a non-zero code
