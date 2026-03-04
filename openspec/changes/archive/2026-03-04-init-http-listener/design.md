## Context

This is the first change to the Ennyn project. There is no existing code — the Go module, entrypoint, and build pipeline all need to be created from scratch. The goal is a minimal HTTP server that will later become the reverse proxy core.

## Goals / Non-Goals

**Goals:**
- Initialize a Go module with a clean project structure
- Create a single-binary HTTP server that listens on a configurable port
- Support `-p` / `-port` CLI flags with a sensible default (7890)
- Bind to `0.0.0.0` by default so the server is reachable from other devices on the network
- Respond with a plain-text "hello world" to all requests
- Shut down gracefully on SIGINT/SIGTERM

**Non-Goals:**
- TLS / HTTPS (later change)
- Reverse proxy routing (later change)
- Configuration file support (later change)
- Custom hostname resolution (later change)

## Decisions

**Go module path: `github.com/ennyn/ennyn`**
Placeholder path that can be updated once a real repo is created. Using a GitHub-style path follows Go conventions.

**Default port: 7890**
Above the well-known range, not used by common dev tools (avoids 3000, 8080, 8443). Easy to remember. Portless uses 1355, Caddy 2019, Traefik 8080 — 7890 is clean and distinctive.

**Standard library only**
`net/http` and `flag` are sufficient. No external dependencies for this change. Keeps the binary small and compilation fast.

**`flag` package for CLI args**
Go's built-in `flag` package supports both `-p 7890` and `-port 7890` by registering two flags pointing to the same variable. Simple and zero-dependency. If richer CLI parsing is needed later, we can switch to `cobra` or `pflag`.

**Graceful shutdown with `os/signal` + `context`**
Listen for SIGINT/SIGTERM, call `http.Server.Shutdown()` with a timeout context. This drains in-flight connections before exit.

## Risks / Trade-offs

**Binding to 0.0.0.0** — Exposes the server to the local network. This is intentional for a dev proxy (e.g., testing from a phone), but worth noting. Non-issue for a development tool; can be restricted to 127.0.0.1 via future config.

**Single-file structure** — Starting with just `main.go` is intentional for this minimal change. Will be restructured into packages (e.g., `cmd/`, `internal/`) as the project grows.
