## Context

Ennyn is a Go development proxy. Today it's a single-file hello-world HTTP server (`main.go`) with no subcommands, no routing, and no process management. We need to turn it into a tool that accepts `ennyn server <host> <app-port> <app-cmd...>`, spawns the app, and reverse-proxies `<host>.localhost` traffic to it.

The existing `http-listener` spec defines the proxy's port binding and graceful shutdown behavior — those stay, but the handler changes from a static response to a host-based reverse proxy.

`*.localhost` domains resolve to `127.0.0.1` per RFC 6761, so no DNS or `/etc/hosts` changes are needed.

## Goals / Non-Goals

**Goals:**
- Single `ennyn server` invocation starts one app behind one `<host>.localhost` hostname
- The proxy reverse-proxies HTTP traffic based on the `Host` header
- The child process is managed (started, monitored, stopped on proxy shutdown)
- No external dependencies beyond the Go standard library

**Non-Goals:**
- Multiple concurrent routes in a single proxy instance (future work)
- HTTPS / TLS termination (future work)
- Configuration files — CLI args only for now
- Daemonization — the proxy runs in the foreground
- Health checking or auto-restarting the child process

## Decisions

### 1. Subcommand model via `flag` subcommands (not cobra/urfave)

Use Go's standard `flag` package with manual subcommand dispatch (`os.Args[1]` switch). The CLI surface is tiny — one subcommand — so a framework adds no value.

**Alternative considered:** cobra — rejected because it's an external dependency for a single subcommand.

### 2. Single reverse proxy per invocation

Each `ennyn server` call manages exactly one host→port mapping. To run multiple apps, the user runs multiple `ennyn server` instances on different proxy ports (or we add multi-route support later).

This keeps the routing logic trivial: if the `Host` header matches `<host>.localhost`, proxy to `127.0.0.1:<app-port>`. Otherwise return 404.

**Alternative considered:** A registry/config file for multiple routes — deferred to keep the first iteration simple.

### 3. `httputil.ReverseProxy` for proxying

Use `net/http/httputil.NewSingleHostReverseProxy` to forward requests. It handles hop-by-hop headers, request rewriting, and error handling out of the box.

### 4. Child process via `os/exec.Cmd`

Start the app command with `exec.Command`, inheriting stdout/stderr from the proxy process so the developer sees app logs inline. On proxy shutdown (SIGINT/SIGTERM), send SIGTERM to the child process group and wait up to 5 seconds before sending SIGKILL.

The command string is split as: first arg is the executable, remaining args are passed through. No shell interpretation — direct exec.

### 5. Host matching includes port stripping

The `Host` header may include a port (e.g., `myapp.localhost:7890`). The router strips the port before comparing against `<host>.localhost`.

## Risks / Trade-offs

- **Child process outlives proxy**: If the proxy crashes hard (SIGKILL), the child may become orphaned. → Mitigation: Use process groups so the OS can clean up. Document this limitation.
- **Port conflicts**: If `<app-port>` isn't actually listening when traffic arrives, the reverse proxy returns 502. → Mitigation: This is expected behavior — the user starts the app command that should bind to that port. No startup health check for now.
- **Single route per instance**: Limits composability. → Mitigation: Acceptable for v1. Multi-route support is a natural follow-up.
