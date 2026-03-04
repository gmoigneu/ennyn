## 1. CLI Subcommand Structure

- [x] 1.1 Restructure main.go to dispatch on `os.Args[1]` subcommand (`server`), with usage error for unknown subcommands
- [x] 1.2 Implement `server` subcommand argument parsing: extract `<host>`, `<app-port>`, `<app-cmd...>` from positional args after flag parsing
- [x] 1.3 Validate host (lowercase alphanumeric + hyphens, no leading/trailing hyphen), port (1–65535), and minimum arg count (>= 3 positional args after `server`)
- [x] 1.4 Support `-p`/`-port` flag on the `server` subcommand for proxy listen port (default 7890)

## 2. Process Manager

- [x] 2.1 Implement child process startup using `os/exec.Cmd` with `Setpgid` for own process group, inheriting stdout/stderr
- [x] 2.2 Implement graceful shutdown: on proxy SIGINT/SIGTERM, send SIGTERM to child process group, wait up to 5s, then SIGKILL
- [x] 2.3 Monitor child process: if child exits unexpectedly, log exit code and shut down the proxy with the same code

## 3. Host-Based Routing

- [x] 3.1 Implement HTTP handler that extracts host from `Host` header (stripping port), compares to `<host>.localhost`
- [x] 3.2 Create `httputil.ReverseProxy` targeting `http://127.0.0.1:<app-port>`, setting `X-Forwarded-For`, `X-Forwarded-Host`, `X-Forwarded-Proto` headers
- [x] 3.3 Return HTTP 404 for non-matching hosts, HTTP 502 custom error handler when upstream is unreachable

## 4. HTTP Listener Changes

- [x] 4.1 Change bind address from `0.0.0.0` to `127.0.0.1`
- [x] 4.2 Replace static "hello world" handler with host-routing handler
- [x] 4.3 Wire everything together: CLI parses args → start child process → start proxy with routing handler → graceful shutdown tears down both

## 5. Testing

- [x] 5.1 Unit tests for argument validation (host format, port range, missing args)
- [x] 5.2 Unit tests for host matching logic (with/without port in Host header, matching, non-matching)
- [x] 5.3 Integration test: start proxy + a simple test server, verify requests to `<host>.localhost` are proxied correctly
