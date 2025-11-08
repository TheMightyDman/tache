Taché

Taché is a fast TUI and CLI for managing dtach sessions on Linux and macOS. Implemented in Go with Bubble Tea. It discovers all dtach sessions owned by your user, not just those created by Taché.

Quick Usage
- `tache` — Open the TUI and browse/attach.
- `tache -l` — List sessions; add `--json` for machine-readable output.
- `tache -a <name-or-prefix>` — Attach by suffix (name) or launch-folder prefix. If multiple match, a TUI filter menu appears.
- `tache start [-n <suffix>] [-- <command...>]` — Create a session (default command: `$SHELL`). If `-n` omitted, we infer from the command or directory.

Discovering All dtach Sessions
- Taché scans for dtach server processes and known socket locations:
  - Process scan (preferred): enumerate dtach servers, extract socket path and launch folder (prefix), verify liveness.
  - Socket scan (fallback): search `~/.dtach`, `/tmp`, and configurable directories. Only show sockets owned by the current user; unverified sockets are hidden by default.
- macOS notes: cwd discovery uses `lsof -p <pid> -a -d cwd`; install with Homebrew if needed.

Suffix (Name) Inference
- If you omit `-n` when starting:
  - Use the first token of the command (e.g., `node`, `bash`, `tail`).
  - Special-case log tails like `tail -f path.log` → `tail:path.log`.
  - Otherwise fall back to the launch directory name.
  - If empty, show as `[none]`; duplicates auto-number, e.g., `[none] (2)`.

More details: see docs/PLAN.md and docs/TROUBLESHOOTING.md (to be expanded).

Build (from source)
- Requires Go 1.22+.
- Build: `go build ./cmd/tache`
- Run TUI: `./tache`
- List: `./tache list --json`
