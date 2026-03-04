## Why

Ennyn needs a foundational HTTP listener before any proxy/routing logic can be built. This is the first step: initialize the Go project and create a minimal HTTP server that binds to a configurable address and port, proving the build pipeline and basic network I/O work.

## What Changes

- Initialize the Go module (`go.mod`)
- Create a `main.go` entrypoint that starts an HTTP server
- Server listens on `0.0.0.0:<port>` by default (port 7890)
- Port is configurable via `-p` or `-port` CLI flag
- Responds to all requests with a plain-text "hello world" message
- Graceful shutdown on SIGINT/SIGTERM

## Capabilities

### New Capabilities
- `http-listener`: Core HTTP server lifecycle — binding to an address/port, accepting connections, and responding to requests. Covers CLI flag parsing for port configuration.

### Modified Capabilities
_None — this is the first change._

## Impact

- New Go module and `go.mod`/`go.sum` files
- New `main.go` entrypoint
- No external dependencies beyond the Go standard library
