## Context

No CI exists. The project builds with `go build` and has no cgo dependencies, making cross-compilation straightforward. The README was written early and doesn't cover the full command set.

## Goals / Non-Goals

**Goals:**
- CI that builds and verifies (`go vet`, `go test`) on every push to main
- Release binaries for 3 platforms on tag push
- Accurate, complete README

**Non-Goals:**
- Automated releases with changelogs (can add later)
- Docker images
- Package manager distribution (homebrew, apt, etc.)

## Decisions

**Single workflow file with build matrix**
One `.github/workflows/build.yml` with a matrix of `GOOS`/`GOARCH` pairs. Simpler than separate workflows per platform. Uses `actions/setup-go` and `actions/upload-artifact`.

**Build targets: linux/amd64, darwin/arm64, windows/amd64**
Linux amd64 covers most servers and WSL. Darwin arm64 covers modern Macs (Apple Silicon). Windows amd64 covers native Windows use. Can add more targets later (darwin/amd64, linux/arm64) if needed.

**Trigger: push to main + tags**
Builds on every push to main for CI verification. On version tags (`v*`), also uploads artifacts. No PR builds to keep it simple.

**Binary naming: `ennyn-<os>-<arch>`**
Consistent, predictable names. Windows binary gets `.exe` suffix.
