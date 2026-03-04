## 1. Instance Registry Package

- [x] 1.1 Create `internal/registry/` package with a function to resolve the run directory (`$XDG_CONFIG_HOME/ennyn/run/` or `~/.config/ennyn/run/`)
- [x] 1.2 Define `Instance` struct: `PID int`, `Host string`, `ProxyPort int`, `AppPort int`, `AppCmd string`, `StartedAt time.Time`
- [x] 1.3 Implement `Register(instance)`: write `<host>.json` to run directory, creating the directory if needed
- [x] 1.4 Implement `Deregister(host)`: remove `<host>.json` from run directory
- [x] 1.5 Implement `IsAlive(pid)`: check if process exists via signal 0
- [x] 1.6 Implement `List()`: read all `.json` files from run directory, check liveness, remove stale files, return live instances
- [x] 1.7 Implement `Get(host)`: load a single instance by hostname, check liveness, clean up if stale
- [x] 1.8 Implement duplicate host check: if state file exists and PID is alive, return an error

## 2. Server Lifecycle Integration

- [x] 2.1 After server starts listening, call `registry.Register()` with the instance metadata
- [x] 2.2 On shutdown (defer), call `registry.Deregister()` to remove the state file
- [x] 2.3 Before starting, check for duplicate host via `registry.Get()` — if alive, exit with error; if stale, clean up and proceed

## 3. List Command

- [x] 3.1 Add `list` case to subcommand dispatch in `main.go`
- [x] 3.2 Implement `runList()`: call `registry.List()`, format and print a table of running instances (host, PID, proxy port, app port, command, uptime)
- [x] 3.3 Print "No running instances" when the list is empty

## 4. Stop Command

- [x] 4.1 Add `stop` case to subcommand dispatch in `main.go`
- [x] 4.2 Implement `runStop(args)`: validate exactly one hostname argument, call `registry.Get(host)`
- [x] 4.3 If instance is alive, send SIGTERM to the PID and print confirmation
- [x] 4.4 If no instance found or stale, print appropriate message and exit accordingly
- [x] 4.5 Update `printUsage()` to include `list` and `stop` commands

## 5. Verification

- [x] 5.1 `go build -o ennyn .` compiles without errors
- [x] 5.2 `go vet ./...` passes
- [x] 5.3 `ennyn server myapp 3000 <cmd>` creates a state file at `~/.config/ennyn/run/myapp.json`
- [x] 5.4 `ennyn list` shows the running instance with correct metadata
- [x] 5.5 `ennyn stop myapp` sends SIGTERM and the instance shuts down, state file removed
- [x] 5.6 `ennyn list` after stop shows no running instances
- [x] 5.7 Starting a duplicate host prints an error and exits non-zero
- [x] 5.8 `ennyn stop nonexistent` prints error and exits non-zero
- [x] 5.9 `ennyn stop` with no args prints usage and exits non-zero
