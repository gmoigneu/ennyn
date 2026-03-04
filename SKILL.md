# Ennyn -- Local Development Proxy

Ennyn is a local development proxy. It gives each of your services a `<name>.localhost` hostname with automatic HTTPS, instead of bare `localhost:<port>` URLs.

## Install

Download the binary for your platform from [GitHub Releases](https://github.com/gmoigneu/ennyn/releases):

- `ennyn-linux-amd64` (Linux / WSL)
- `ennyn-darwin-arm64` (macOS Apple Silicon)
- `ennyn-windows-amd64.exe` (Windows)

Place it on your PATH. No runtime dependencies required.

## One-time setup

Run `ennyn trust` to generate a local CA and install it in the system trust store. This requires sudo/admin. It creates a wildcard certificate for `*.localhost` so all services get trusted HTTPS automatically.

## Usage

### Config file (recommended)

Create `ennyn.yml` in the project root:

```yaml
services:
  - host: web
    port: 3000
    command: npm run dev
  - host: api
    port: 4000
    command: go run ./cmd/api
```

Then run:

```
ennyn serve
```

This starts all services and proxies them:
- `web.localhost:7890` -> `127.0.0.1:3000`
- `api.localhost:7890` -> `127.0.0.1:4000`

### Single service

```
ennyn server myapp 3000 npm run dev
```

Proxies `myapp.localhost:7890` -> `127.0.0.1:3000`.

### Management

```
ennyn list              # Show running instances
ennyn stop <host>       # Stop a specific service
ennyn stop              # Stop all services from ennyn.yml
```

## Config format

```yaml
services:
  - host: <name>        # Required. Lowercase alphanumeric + hyphens. Becomes <name>.localhost
    port: <int>         # Required. App port (1-65535)
    command: <string>   # Required. Shell command to start the service
```

Hosts must be unique. The proxy listens on port 7890 by default (override with `-p`).

## File locations

- CA and certs: `~/.config/ennyn/ca/`
- Instance state: `~/.config/ennyn/run/`

## Commands reference

| Command | Description |
|---------|-------------|
| `ennyn serve [-p port]` | Start all services from `ennyn.yml` |
| `ennyn server [-p port] <host> <app-port> <cmd> [args...]` | Start a single service |
| `ennyn list` | List running instances |
| `ennyn stop [<host>]` | Stop a service or all services |
| `ennyn trust` | Install local CA and certificates |
| `ennyn trust --uninstall` | Remove CA and certificates |
