## 1. GitHub Actions

- [x] 1.1 Create `.github/workflows/build.yml` with a build matrix for linux/amd64, darwin/arm64, windows/amd64
- [x] 1.2 Workflow triggers on push to main and on tags matching `v*`
- [x] 1.3 Steps: checkout, setup-go, `go vet ./...`, `go test ./...`, cross-compile with `GOOS`/`GOARCH`, upload artifacts
- [x] 1.4 Binary naming: `ennyn-<os>-<arch>` (`.exe` suffix for windows)

## 2. README Rewrite

- [x] 2.1 Update overview: remove "planned capabilities" framing, describe what actually works now (proxy routing, TLS, list/stop)
- [x] 2.2 Fix prerequisites: Go is only needed for building from source / contributing, not for running
- [x] 2.3 Add all current commands to the Commands table: `server`, `list`, `stop`, `trust`, `trust --uninstall`
- [x] 2.4 Document the `-p`/`-port` flag for server
- [x] 2.5 Add `list` and `stop` usage examples
- [x] 2.6 Update the "Stop" section to mention both `ennyn stop <host>` and Ctrl+C
- [x] 2.7 Keep the cross-compilation section, add windows target

## 3. Housekeeping

- [x] 3.1 Add `.gitignore` with the `ennyn` binary and common Go build artifacts
