## Why

When running multiple services, developers have no way to see which Ennyn instances are active or to stop a specific one by hostname. They must manually find and kill processes. `ennyn list` and `ennyn stop <host>` solve this with a simple PID-file registry.

## What Changes

- Each `ennyn server` instance writes a state file on startup and removes it on shutdown
- Add `ennyn list` command that reads state files and shows running instances (host, PID, ports, status)
- Add `ennyn stop <host>` command that sends SIGTERM to the instance managing that host
- State files stored at `~/.config/ennyn/run/<host>.json` (respects `$XDG_CONFIG_HOME`)
- Stale state files (process no longer running) are detected and cleaned up

## Capabilities

### New Capabilities
- `instance-registry`: State file management for tracking running Ennyn server instances. Covers writing, reading, cleanup, and the state directory convention.
- `list-stop-commands`: The `ennyn list` and `ennyn stop <host>` CLI commands.

### Modified Capabilities
- `http-listener`: Server startup must now write a state file, and shutdown must remove it.

## Impact

- New `internal/registry/` package for state file operations
- New subcommands in `main.go` (`list`, `stop`)
- Modified server lifecycle in `main.go` (register on start, deregister on shutdown)
- New state directory at `~/.config/ennyn/run/`
