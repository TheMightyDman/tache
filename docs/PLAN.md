Taché — Project Plan (Go)

Overview
- Taché is a terminal session manager for dtach. It lists and manages active dtach sessions and provides a fast TUI (Bubble Tea-based) to browse, filter, attach, create, rename, and kill sessions.
- Implemented in Go for small, fast, single binaries. Targets Linux and macOS.

Performance & Refresh
- TUI refresh is debounced (75–150 ms) when typing/filtering to keep UI responsive.
- Discovery uses incremental scanning:
  - Cache previous results in `index.json` with `seenAt` timestamps.
  - On refresh, re-check only changed PIDs and new dtach processes; fully rescan on demand.
  - macOS `lsof -p <pid> -a -d cwd` results are cached per-PID with a short TTL (e.g., 10s) to avoid repeated calls.
  - Linux uses `/proc/<pid>/cwd` reads, which are cheap; still memoized per-PID during a refresh.

CLI
- `tache` — Open TUI listing all active sessions.
- `tache -l [--json]` — List sessions to stdout. JSON for scripting.
- `tache -a <prefix-or-suffix>` — Attach by suffix or by launch-folder prefix. If multiple matches, open TUI pre-filtered. If single match, attach directly.
- `tache start [-n <suffix>] [-C <dir>] [-- <command...>]` — Create new session (default command: `$SHELL`).
- `tache rename <id|suffix> <new-suffix>` — Rename the metadata label only (does not rename the socket) to avoid service disruption.
- `tache kill <id|suffix>` — Kill the session’s underlying process (with confirmation; supports `--signal`).
- Global flags: `--dtach-bin`, `--search-dir`, `--state-dir`, `--detach-char`, `--no-color`, `--verbose`, `--debug`.

Naming Rules (Suffix Inference)
- If `-n <suffix>` not provided for `tache start`:
  - If a command is provided, use the start of the command (first token) as the suffix (e.g., `node`, `bash`, `tail`).
  - If the session is a console/log tail (e.g., `tail -f <file>`), infer a descriptive suffix such as `tail:<file-basename>` when feasible.
  - Otherwise infer from launch directory name.
  - As a fallback, allow empty suffix. If empty suffix already exists, auto-number them (e.g., `", ", (2)`, etc.).
- Suffix collisions are allowed; disambiguate by prefix (launch folder) and socket path.

Session Model
- Identity: socket path (authoritative), dtach server PID, child command, createdAt, updatedAt.
- Metadata: prefix (launch folder), suffix (friendly label), detach char, lastAttachAt, notes.
- Storage: XDG state/config locations.
  - Linux: `~/.local/state/tache/index.json` and `~/.config/tache/config.json`.
  - macOS: `~/Library/Application Support/tache/state/index.json` and `~/Library/Application Support/tache/config.json`.

Discovery & Scanning
- Discover all dtach sessions for the current user (not only Taché-created sessions).
- Primary method: process scan for dtach server processes.
  - Linux: `ps` + `/proc/<pid>/cwd` to obtain the launch folder (prefix).
  - macOS: `ps` + `lsof -p <pid> -a -d cwd` to obtain the launch folder (prefix).
  - Parse dtach args to locate socket path and determine mode.
  - Confirm health by checking socket existence and server PID liveness.
- Secondary method: socket directory scanning.
  - Robust search of common locations: `~/.dtach`, `/tmp`, `~/.cache/dtach`, and user-configured paths.
  - Validate candidate sockets WITHOUT writing to the socket (do not use `dtach -p`, which writes stdin to the session):
    - Check the path is a UNIX socket (`stat`), is owned by current user, and has reasonable permissions.
    - Cross-reference with process scan: prefer sockets that are referenced by a running dtach server.
    - If unmatched, mark as `unverified` and hide by default (toggle to show).
- Stale handling: mark sessions stale when PID is dead or socket missing; hidden by default, toggleable.
- Performance: cache scan results in `index.json` with timestamps; incremental refresh; configurable scan intervals.

Attach Behavior
- TUI attach: spawn `dtach -a <socket>` as a child attached to the same TTY; on detach or command exit, return to TUI.
- CLI `-a`: resolve suffix/prefix to sessions. If single match, attach directly; if multiple, open TUI pre-filtered.
- Environment inside attached sessions: set `TACHE_SOCKET`, `TACHE_SESSION_ID`, `TACHE_FROM_TUI=1`.
- Detach char: default to `^\` (Ctrl-\), configurable per session/attach.

Create / Rename / Kill
- Create (`tache start`):
  - Defaults: cwd=`$PWD`, command=`$SHELL`, suffix inferred per rules above.
  - Socket path defaults under `~/.local/state/tache/sockets/<slug>.sock` (Linux) and corresponding macOS state dir to avoid clutter.
  - Launch with `dtach -n <socket> -- <command...>` and record metadata into index.
- Rename (`tache rename`):
  - Update only Taché metadata (suffix) to avoid disrupting the service. Socket rename is not performed by default.
  - If session is killed, delete its metadata automatically.
- Kill (`tache kill`):
  - Send signal to underlying process (default SIGTERM), with confirmation UI.

Concurrency & Integrity
- All writes to `index.json` are atomic: write to `index.json.tmp` then `rename()`.
- Use an advisory lock file `index.lock` (flock) around write operations; retry with jitter and stale lock detection.
- Concurrent Taché instances will serialize writes; readers tolerate absence of lock.

Pruning Stale Sessions
- Command: `tache prune [--all | --older-than <days>]` removes stale entries from `index.json` and cleans up orphaned tache-managed socket files.
- Auto-prune: optional config to auto-remove stale entries older than N days on startup.

TUI (Go Bubble Tea)
- Libraries: `bubbletea` (core), `bubbles` (list, textinput), `lipgloss` (styling).
- Layout:
  - Header: Taché title and quick help.
  - Filter bar: fuzzy search across suffix, prefix, command.
  - List: [Status, Suffix, Prefix, PID, Uptime, Socket].
  - Footer: keybinds and status.
- Keybinds:
  - Enter/a: attach
  - f: focus filter; Esc clears
  - c: create session modal (name, cwd, command, detach char)
  - r: rename suffix
  - k: kill (confirm)
  - s: sort by suffix/prefix/uptime
  - j/k or arrows: navigate; g/G: page up/down; u: refresh; ?: help overlay
- Behavior:
  - Live refresh interval + filesystem events where available.
  - Badges for stale/detached/busy/new.
  - Disambiguation: collisions shown as `suffix (~/relative/path)`; Empty suffix rendered as `[none]`, then `[none] (2)`, `[none] (3)`.

Loop Avoidance & Interop
- Prevent infinite loops when run inside Taché or tmux:
  - Detect `TACHE_FROM_TUI` and warn before nesting; return to menu after detach automatically by spawning attach as a child.
  - Detect `TMUX`; provide optional `tache switch` to detach and open menu in same pane.
- Optional helper: `tache back` — when in a Taché session, detach and re-open the menu (best-effort using env/socket info).
 - General nesting: maintain `TACHE_NESTING_LEVEL` env and warn when >1, also detect `STY` (screen) and `ZELLIJ` to inform the user.

Configuration
- File: `~/.config/tache/config.json` (Linux) or `~/Library/Application Support/tache/config.json` (macOS).
- Keys: dtach path, search paths, default socket dir, detach char, TUI defaults (sort, columns), refresh interval.
- Env overrides: `TACHE_CONFIG`, `TACHE_STATE_DIR`, `TACHE_SEARCH_DIRS`, `TACHE_DTACH_BIN`.

Packaging & Installation
- Produce static single binaries via Go builds; distribute Linux x64/arm64 and macOS x64/arm64.
- Use GoReleaser for CI packaging, checksums, and Homebrew tap.

Testing Strategy
- Unit: arg parsing, suffix inference, filtering, state indexing.
- Integration: spawn real dtach sessions on CI (Linux runner); validate discovery, attach/detach, create/rename/kill.
- macOS: rely on community runners to validate `lsof` cwd discovery and behavior.
- TUI smoke: render non-interactive snapshots and keymap tests.

Timeline (MVP-first)
- Milestone 1: CLI skeleton (`-l`, `--json`, discovery, index store, suffix inference).
- Milestone 2: TUI list + attach flow (OpenTUI), filters/sort, return-to-menu.
- Milestone 3: create/rename/kill, `-a` logic, config and env overrides.
- Milestone 4: robust scanning and stale handling, macOS cwd via `lsof`, polish.
- Milestone 5: packaging and prebuilt binaries, docs.

Docs To Add (README and Usage)
- Explain discovering all dtach sessions (even those not created via Taché) and how to configure search directories.
- Examples:
  - `tache -l --json` and parsing output.
  - `tache -a api` to attach by suffix, and handling multiples.
  - `tache start -- node server.js` inferring suffix `node`.
  - `tache start -n logs -- tail -f app.log` inferring `tail:app.log` if `-n` omitted.
- Notes on macOS requirements: `lsof` needed for cwd discovery; Homebrew provides it.

Failure Modes & Recovery
- `index.json` corrupted: keep a rotating backup `index.json.bak`; on parse failure, move the bad file to `index.json.bad.<ts>`, rebuild from discovery, and continue in read-only mode until next write.
- `dtach` not found: clearly report and show installation guidance (Linux/macOS). Allow read-only listing from cache where possible.
- Session dies during attach: on child exit or attach failure, return to TUI with a clear status and offer to prune.

Schema Versioning
- `index.json` includes a `schemaVersion` field. On startup, migrate older versions in-memory and write back using the current schema.
- Maintain simple, forward-only migrations with tests; create backups before migration.

TUI Dependency
- TUI is implemented with Bubble Tea. CLI remains fully functional without TUI (list/attach/start/kill). No alternate TUI framework planned.

Detach Character Conflicts
- Default detach is `^\` (Ctrl-\). During creation, allow per-session override.
- Warn if user selects a high-conflict char (e.g., `^C`, `^D`) and require confirmation.

macOS CI
- Add GitHub Actions macOS job to run discovery and CLI tests. If `dtach` is unavailable via package manager, build from source in CI or gate attach tests while still testing discovery and indexing.

Socket Path Guidance
- Recommend tache-managed sockets under the state dir for persistence and isolation.
- Document pros/cons: state dir survives reboots; `/tmp` is ephemeral and cleaned on reboot; `~/.dtach` is conventional but shared.

Docs & Troubleshooting
- Add docs covering migration from raw dtach usage, common errors (sessions not visible, attach failures, permission issues), and shell integration examples (aliases, auto-naming templates).

Filtering at Scale
- Fuzzy search uses a lightweight scorer with pre-tokenized fields (suffix, prefix, command). Debounced input and list virtualization keep 100–1000 sessions responsive.

Bulk Operations
- CLI: support `tache kill --filter '<pattern>' --yes` to operate on multiple matches. Future: tags and multi-select in TUI.

Open Questions (tracked for future)
- macOS socket scanning defaults (common paths to add beyond `$TMPDIR`?): gather community input.
- Optional socket rename for Taché-managed sessions behind a flag, with downtime warning.
