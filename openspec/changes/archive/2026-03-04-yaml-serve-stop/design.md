## Context

Ennyn currently requires a separate `ennyn server <host> <port> <cmd>` invocation per service. For multi-service projects this means manually running several commands. There is no declarative project-level configuration. The `ennyn stop` command requires an explicit hostname argument.

The instance registry (`~/.config/ennyn/run/<host>.json`) already tracks running instances, and `ennyn list`/`ennyn stop <host>` already manage them. This change adds a config layer on top.

## Goals / Non-Goals

**Goals:**
- Define a simple, flat YAML config file (`ennyn.yml`) listing services
- `ennyn serve` starts all services from the config as background processes
- `ennyn stop` (no args) stops all services defined in the config
- Reuse existing `ennyn server` machinery — the serve command orchestrates server instances

**Non-Goals:**
- Hot-reload / file-watching on config changes
- Dependency ordering between services
- Health checks or readiness probes
- Merging multiple config files or config inheritance
- `ennyn.yaml` as alternate filename (only `ennyn.yml`)

## Decisions

### 1. Config file format: flat service list

```yaml
services:
  - host: myapp
    port: 3000
    command: npm run dev
  - host: api
    port: 8080
    command: go run ./cmd/api
```

**Rationale:** Minimal surface area. Each service maps directly to `ennyn server` arguments. The `services` top-level key allows future extension (e.g., global settings) without breaking changes.

**Alternatives considered:**
- Map keyed by host name (`myapp: {port: 3000, ...}`) — slightly more compact but prevents duplicate-host validation at parse time and makes ordering ambiguous.
- TOML — Go has good TOML support, but YAML is more familiar to the target audience (web developers).

### 2. `ennyn serve` starts background server processes

`ennyn serve` will fork `ennyn server` instances as background processes (using `os/exec` with `Cmd.Start()`), one per service entry. It will not manage the processes itself — the existing per-instance process manager and instance registry handle lifecycle.

**Rationale:** Reuses all existing infrastructure. Each server instance writes its own state file, responds to signals, and manages its child process. `ennyn serve` is purely an orchestrator.

**Alternatives considered:**
- Single-process multi-service mode — would require significant refactoring of the process manager and HTTP listener. Much more complex for marginal benefit.

### 3. `ennyn stop` without args reads config

When `ennyn stop` receives no arguments, it reads `ennyn.yml`, extracts all hostnames, and calls the existing stop logic for each one. This is additive — `ennyn stop <host>` continues to work unchanged.

**Rationale:** Symmetric with `ennyn serve`. Natural user expectation: "serve starts everything, stop stops everything."

### 4. YAML parsing with `gopkg.in/yaml.v3`

**Rationale:** Mature, widely used, no cgo. Standard choice for Go YAML parsing.

### 5. Config file discovery: current directory only

`ennyn serve` looks for `ennyn.yml` in the working directory. No upward directory traversal.

**Rationale:** Simple, predictable. Matches behavior of tools like `docker-compose.yml`. Users run commands from the project root.

## Risks / Trade-offs

- **Naming confusion: `serve` vs `server`** — The two commands are similar in name. `serve` is the multi-service config command; `server` is the single-service direct command. Mitigation: Clear help text distinguishing the two.
- **Background process visibility** — `ennyn serve` starts processes and returns. Users must use `ennyn list` to see status. Mitigation: Print a summary table after starting services.
- **Partial start failure** — If one service fails to start (e.g., port conflict), others may already be running. Mitigation: Start all, report errors, let user fix and re-run. Do not roll back successful starts.
