## Why

Ennyn has no CI pipeline and the README is outdated — it's missing the `list` and `stop` commands, incorrectly lists Go as a prerequisite for users (it's only needed for contributing), and doesn't reflect the current feature set. We need GitHub Actions to build release binaries for all 3 target platforms and a README that accurately documents what exists.

## What Changes

- Add GitHub Actions workflow that builds binaries for linux/amd64, darwin/arm64, and windows/amd64 on push to main and on tags
- Rewrite `docs/README.md` to reflect all current commands (`server`, `list`, `stop`, `trust`), correct prerequisites (no Go needed to run — it's a single binary), and the `-p`/`-port` flag
- Add a `.gitignore` for the built binary

## Capabilities

### New Capabilities
_None — this is a CI/docs change with no new runtime capabilities._

### Modified Capabilities
_None._

## Impact

- New `.github/workflows/build.yml`
- Updated `docs/README.md`
- New `.gitignore`
