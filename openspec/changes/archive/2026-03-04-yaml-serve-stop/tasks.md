## 1. Config Parser

- [x] 1.1 Add `gopkg.in/yaml.v3` dependency
- [x] 1.2 Create `internal/config/config.go` with `Service` struct and `Load(path string)` function that reads and parses `ennyn.yml`
- [x] 1.3 Implement validation: host format, port range, non-empty command, duplicate host detection
- [x] 1.4 Implement command string splitting (split on whitespace into executable + args)
- [x] 1.5 Write tests for config parsing: valid config, empty services, invalid host, invalid port, missing command, duplicate hosts, missing file

## 2. Serve Command

- [x] 2.1 Create `internal/serve/serve.go` with `Run(configPath string, proxyPort int)` function
- [x] 2.2 Implement service startup: for each config service, exec `ennyn server` as a detached background process (using `os/exec.Cmd` with `SysProcAttr` for detach, no stdin)
- [x] 2.3 Skip already-running services by checking instance registry before starting
- [x] 2.4 Print summary table after starting (hostname, port, status)
- [x] 2.5 Wire `ennyn serve` subcommand in `main.go` with `-p`/`-port` flag support
- [x] 2.6 Write tests for serve logic: config loading, skip-already-running detection

## 3. Stop Config Mode

- [x] 3.1 Modify `runStop` in `main.go`: when no args provided, attempt to load `ennyn.yml` from current directory
- [x] 3.2 If config found, iterate over services and stop each one (reuse existing stop logic), printing per-service results
- [x] 3.3 If no config found and no args, print usage and exit non-zero (current behavior)
- [x] 3.4 Write tests for stop-all-from-config behavior

## 4. Integration

- [x] 4.1 Update `printUsage()` in `main.go` to include `ennyn serve` and updated `ennyn stop` usage
- [x] 4.2 Manual end-to-end test: create `ennyn.yml`, run `ennyn serve`, verify with `ennyn list`, run `ennyn stop`
