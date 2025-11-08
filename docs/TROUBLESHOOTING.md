Troubleshooting

Sessions Don’t Appear
- Ensure `dtach` is installed and on PATH.
- On macOS, install `lsof` (Homebrew) for accurate cwd discovery.
- Check that your sockets are owned by your user; Taché ignores sockets owned by others.
- If sockets are in nonstandard locations, add them to `searchDirs` in config.

Cannot Attach
- The session may be stale (dead PID). Use `tache -l --json` to confirm status.
- Try `tache prune` to remove stale entries.
- Verify socket path exists and is accessible; if on `/tmp`, it may have been cleaned after reboot.

Corrupted State
- If `index.json` is corrupted, Taché will rebuild from discovery and keep a backup. You can remove the `.bad.*` file after verifying.

Nested UIs
- If you launched Taché from inside a Taché-attached session, you’ll see a nesting warning. Prefer returning to the menu after detach or use `tache back`.
