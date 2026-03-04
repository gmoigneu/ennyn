# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Ennyn ("Doors" in Sindarin) is an open-source development proxy that runs locally on a developer's machine. It manages hostname resolution, TLS certificates, and request routing for applications composed of multiple services running on different ports. Target platforms: macOS and Windows (WSL).

## Language & Build

- **Language:** Go
- **Module path:** TBD (set in go.mod once initialized)
- **Build:** `go build ./...`
- **Test all:** `go test ./...`
- **Test single:** `go test ./path/to/package -run TestName`
- **Lint:** `golangci-lint run` (install via https://golangci-lint.run)
- **Vet:** `go vet ./...`
- **Cross-compile:** `GOOS=darwin GOARCH=arm64 go build -o ennyn-darwin-arm64` / `GOOS=linux GOARCH=amd64 go build -o ennyn-linux-amd64`

## Architecture

*To be updated as the codebase grows.* Key subsystems planned:

- **Proxy/Router:** Reverse proxy that routes incoming requests to local services based on hostname or path rules
- **Certificate Management:** Generates and manages local TLS certificates (likely a local CA) so services can be accessed over HTTPS
- **Host Management:** Configures local DNS/hostname resolution (e.g., /etc/hosts or platform-specific resolver) to point custom domains at the proxy
- **Configuration:** Declarative config (format TBD) defining routes, hostnames, and upstream service addresses

## Platform Considerations

- Must compile and run on macOS (native) and Linux (for WSL on Windows)
- Host file and certificate trust store manipulation differs per OS — use build tags or runtime detection as needed
- Avoid cgo dependencies where possible to simplify cross-compilation
