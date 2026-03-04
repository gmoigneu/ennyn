## Why

Running `ennyn server` once per service is tedious for multi-service projects. Developers need a single command to spin up all services defined in a project config file, and a matching command to tear them all down.

## What Changes

- Introduce `ennyn.yml` project configuration file format defining services (host, port, command)
- Add `ennyn serve` command that reads `ennyn.yml` from the current directory and starts all listed services as background `ennyn server` instances
- Extend `ennyn stop` to accept no arguments — when invoked without a host, it reads `ennyn.yml` and stops all services defined in it

## Capabilities

### New Capabilities
- `project-config`: Defines the `ennyn.yml` file format — a list of services each specifying host, port, and command
- `serve-command`: The `ennyn serve` CLI command that reads `ennyn.yml` and starts all services

### Modified Capabilities
- `list-stop-commands`: `ennyn stop` gains a no-argument mode that reads `ennyn.yml` and stops all services defined in it

## Impact

- New YAML dependency (`gopkg.in/yaml.v3` or `sigs.k8s.io/yaml`)
- New CLI subcommand `serve` alongside existing `server`
- Behavioral change to `ennyn stop`: no-argument invocation now valid (currently requires a host argument)
- No breaking changes to existing commands — `ennyn server`, `ennyn stop <host>`, and `ennyn list` continue to work unchanged
