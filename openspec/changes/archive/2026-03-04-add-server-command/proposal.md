## Why

Ennyn currently has no way to manage application processes or route traffic to them. Developers need a single command to launch an app behind a `.localhost` hostname so they can access it as `myapp.localhost` instead of `localhost:PORT`. This is the core feature that makes Ennyn useful as a development proxy.

## What Changes

- Add `ennyn server <host> <app-port> <app-cmd...>` CLI command
- The command starts `<app-cmd>` as a background child process
- The proxy routes HTTP requests arriving at `<host>.localhost` to `127.0.0.1:<app-port>`
- Host-based routing using the `Host` header to multiplex multiple apps on a single proxy port
- The proxy manages the lifecycle of the child process (stop it on shutdown)

## Capabilities

### New Capabilities
- `server-command`: CLI parsing for `ennyn server <host> <app-port> <app-cmd...>`, argument validation, and command dispatch
- `process-manager`: Starting the app command as a child process, forwarding signals, and stopping it on proxy shutdown
- `host-routing`: Reverse-proxy routing based on `Host` header matching `<host>.localhost`, forwarding requests to `127.0.0.1:<app-port>`

### Modified Capabilities
- `http-listener`: The listener must now use the host-routing handler instead of the static "hello world" response. The proxy port becomes the entry point for all routed traffic.

## Impact

- **CLI**: Moves from a bare `ennyn` invocation to a subcommand model (`ennyn server ...`)
- **main.go**: Will be restructured to support subcommands
- **Dependencies**: `os/exec` for process management; the existing `net/http` reverse proxy utilities (`httputil.ReverseProxy`) for routing — no new external dependencies expected
- **Networking**: Relies on `*.localhost` resolving to `127.0.0.1`, which is guaranteed by RFC 6761 and works out of the box on modern OSes — no `/etc/hosts` manipulation needed
