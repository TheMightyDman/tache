package discovery

import (
    "context"
    "errors"
    "fmt"
    "os/exec"
    "strings"
    "time"

    "tache/internal/types"
)

// Discover returns currently known dtach sessions. This is a scaffold implementation
// that returns an empty list. It will be expanded to scan processes and sockets.
func Discover(ctx context.Context) ([]types.Session, error) {
    // TODO: implement process scan and socket validation.
    return []types.Session{}, nil
}

// AttachBySelector attaches to a matching session (by suffix or prefix). Stub.
func AttachBySelector(selector string) error {
    // TODO: implement matching and attach using dtach -a <socket>
    return errors.New("attach not implemented yet")
}

// StartSession creates a new dtach session. Stub.
func StartSession(chdir string, name string, command []string) error {
    // Example of how dtach create might look, to be expanded with naming rules.
    // dtach -n <socket> -- <command...>
    _ = chdir
    _ = name
    _ = command
    return errors.New("start not implemented yet")
}

// Rename updates metadata for a session. Stub.
func Rename(selector string, newName string) error {
    _ = selector
    _ = newName
    return errors.New("rename not implemented yet")
}

// Kill terminates a session's process. Stub.
func Kill(selector string, yes bool) error {
    _ = selector
    _ = yes
    return errors.New("kill not implemented yet")
}

// Prune removes stale sessions/metadata. Stub.
func Prune(all bool, olderThanDays int) error {
    _ = all
    _ = olderThanDays
    return errors.New("prune not implemented yet")
}

// helper to run external commands with timeout
func runWithTimeout(timeout time.Duration, name string, args ...string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    cmd := exec.CommandContext(ctx, name, args...)
    out, err := cmd.CombinedOutput()
    return strings.TrimSpace(string(out)), err
}

