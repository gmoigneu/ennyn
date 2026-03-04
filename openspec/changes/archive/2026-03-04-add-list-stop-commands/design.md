## Context

Ennyn currently runs one server instance per invocation with no shared state. Each `ennyn server` is independent — there's no way to discover or control running instances. The codebase uses `~/.config/ennyn/` (via `$XDG_CONFIG_HOME`) for CA storage, so we extend this convention for runtime state.

## Goals / Non-Goals

**Goals:**
- Track running instances via PID files in a state directory
- Provide `ennyn list` to show all active instances
- Provide `ennyn stop <host>` to gracefully terminate a specific instance
- Detect and clean up stale state files (crashed processes)

**Non-Goals:**
- Daemon mode or central coordinator
- Inter-process communication between instances
- Stopping all instances at once (can be added later)
- Watching/auto-restart of stopped instances

## Decisions

**PID files as JSON state files**
Each running instance writes `~/.config/ennyn/run/<host>.json` containing PID, hostname, proxy port, app port, app command, and start time. JSON is easy to parse and allows rich metadata beyond just a PID. Files are named by hostname for easy lookup.

**State directory: `~/.config/ennyn/run/`**
Follows the existing XDG convention used by the CA storage. `run/` is a common name for runtime state. Respects `$XDG_CONFIG_HOME`.

**Stale detection via `os.FindProcess` + signal 0**
On Unix, sending signal 0 to a PID checks if the process exists without affecting it. If the PID in the state file is not running, the file is stale and should be cleaned up. This is checked during `list` and `stop`.

**`stop` sends SIGTERM, not SIGKILL**
The server already handles SIGTERM gracefully (drains connections, stops child process). SIGTERM is the right signal. If the process doesn't exit within a reasonable time, the user can escalate manually.

**Deregistration via `defer` in server lifecycle**
The state file is written immediately after the server starts listening, and removed via `defer` in the server function. This covers normal exit, signal-based shutdown, and most crash scenarios. If the process is killed with SIGKILL, the stale file is cleaned up on next `list` or `stop`.

**File locking: not needed**
Each hostname maps to exactly one state file, and only one instance should run per hostname. If a state file already exists and the PID is alive, `ennyn server` should refuse to start (duplicate host detection). This prevents conflicts without locking.

## Risks / Trade-offs

**SIGKILL or power loss leaves stale files** → Detected and cleaned up by `list` and `stop` via process liveness checks. Acceptable for a dev tool.

**Race condition on startup** → Two instances starting the same hostname simultaneously could both check "no file exists" then both write. Extremely unlikely in practice for a local dev tool. File locking could be added later if needed.

**PID reuse** → After a crash, the OS could reassign the PID to an unrelated process. Mitigated by storing the start time in the state file and comparing with process start time where possible. For a dev tool, this edge case is acceptable.
