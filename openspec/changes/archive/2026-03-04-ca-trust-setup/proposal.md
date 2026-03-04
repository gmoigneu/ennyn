## Why

Ennyn needs to serve HTTPS to local services with certificates browsers actually trust. This requires a local Certificate Authority whose root certificate is installed in the OS trust store. Without this, browsers show security warnings and service-to-service calls fail TLS verification.

## What Changes

- Add an `ennyn trust` command that generates a local root CA (ECDSA P-256) and installs it into the platform trust store
- Also generate a wildcard leaf certificate for `*.localhost` (plus `localhost`, `127.0.0.1`, `::1`) signed by the CA — a single cert covers all `<app>.localhost` hostnames
- CA and leaf cert stored at `~/.config/ennyn/` (respects `$XDG_CONFIG_HOME`)
- Supports macOS (System Keychain), Linux (ca-certificates), and Windows (system cert store via crypt32)
- Auto-detects WSL and offers to install the CA in the Windows trust store as well
- Add `ennyn trust --uninstall` to remove the CA from trust stores and delete the local files
- After setup, print clear instructions about what was done and next steps
- New external dependency: `github.com/smallstep/truststore` for platform trust store operations
- No hostname resolution needed: `.localhost` subdomains resolve to `127.0.0.1` automatically per RFC 6761

## Capabilities

### New Capabilities
- `ca-trust`: Local CA generation, wildcard leaf cert for `*.localhost`, trust store installation/uninstallation, WSL detection, and storage management. Covers the `ennyn trust` CLI command.

### Modified Capabilities
_None._

## Impact

- New `internal/ca/` package for CA generation and storage
- New `cmd/trust.go` (or equivalent) for the CLI command
- Refactor `main.go` to use subcommands (currently only has the server; now needs `trust` alongside the server start)
- New dependency: `smallstep/truststore`
- Requires sudo/admin privileges when run (trust store writes are privileged)
