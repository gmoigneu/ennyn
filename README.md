# Ennyn

A local development proxy that routes requests to your services by hostname, with automatic HTTPS.

## Overview

Ennyn ("Doors" in Sindarin) runs on your development machine and sits in front of your local services. Instead of juggling `localhost` ports, you give each service a hostname like `myapp.localhost` and access it over HTTPS.

All services use `<name>.localhost` hostnames, which resolve to `127.0.0.1` automatically on all platforms (per RFC 6761). No `/etc/hosts` editing required.

Features:

- **Reverse proxy** -- route requests to local services by hostname
- **Automatic HTTPS** -- local CA with a wildcard `*.localhost` certificate
- **Process management** -- starts your app and proxies to it in one command
- **Instance tracking** -- list running instances, stop them by name
- **Cross-platform** -- single binary for macOS, Linux, and Windows

## Installation

Download the latest binary for your platform from [Releases](https://github.com/gmoigneu/ennyn/releases), or build from source (see [Contributing](#contributing)).

## Getting Started

### 1. Set up HTTPS (one-time)

```
ennyn trust
```

Generates a local Certificate Authority and a wildcard certificate for `*.localhost`, then installs the CA into your system trust store. Requires sudo/admin access.

On WSL, Ennyn detects the environment and offers to install the CA in the Windows trust store as well.

To remove the CA and all certificates:

```
ennyn trust --uninstall
```

### 2. Start a service

```
ennyn server myapp 3000 npm run dev
```

This starts `npm run dev`, then proxies `myapp.localhost:7890` to `127.0.0.1:3000`. Run multiple services in separate terminals:

```
ennyn server api 4000 go run ./cmd/api
ennyn server web 3000 npm run dev -p 3001
```

### 3. Manage running instances

```
ennyn list
```

```
HOST  PID    PROXY PORT  APP PORT  COMMAND  UPTIME
api   12345  7890        4000      go       2m30s
web   12346  3001        3000      npm      1m15s
```

Stop a specific instance:

```
ennyn stop api
```

Or press `Ctrl+C` in the terminal where the instance is running.

## Commands

| Command | Description |
|---------|-------------|
| `ennyn server [flags] <host> <app-port> <app-cmd> [args...]` | Start the proxy and app |
| `ennyn list` | Show all running instances |
| `ennyn stop <host>` | Stop a running instance by hostname |
| `ennyn trust` | Generate CA and certificates, install into trust store |
| `ennyn trust --uninstall` | Remove CA from trust store and delete certificate files |

### Server flags

| Flag | Default | Description |
|------|---------|-------------|
| `-p`, `-port` | `7890` | Proxy port to listen on |

## File locations

All files stored under `$XDG_CONFIG_HOME/ennyn/` (defaults to `~/.config/ennyn/`):

| Path | Description |
|------|-------------|
| `ca/ca.crt` | Root CA certificate |
| `ca/ca.key` | Root CA private key |
| `ca/localhost.crt` | Wildcard cert for `*.localhost` |
| `ca/localhost.key` | Wildcard cert private key |
| `run/<host>.json` | State file for running instance |

## Contributing

### Prerequisites

- [Go](https://go.dev/dl/) 1.25 or later

### Build from source

```
git clone https://github.com/gmoigneu/ennyn.git
cd ennyn
go build -o ennyn .
```

### Cross-compilation

Ennyn compiles to a single binary with no cgo dependencies:

```
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o ennyn-darwin-arm64 .

# Linux (x86_64, for WSL)
GOOS=linux GOARCH=amd64 go build -o ennyn-linux-amd64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o ennyn-windows-amd64.exe .
```

### Run tests

```
go test ./...
go vet ./...
```
