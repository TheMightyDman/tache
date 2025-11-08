package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "tache/internal/discovery"
    "tache/internal/tui"
)

var (
    flagJSON  bool
    rootList  bool
    rootJSON  bool
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "tache",
        Short: "Taché — dtach session manager",
        RunE: func(cmd *cobra.Command, args []string) error {
            // Root flags behavior: -l/--list and --json
            if rootList {
                ctx := context.Background()
                sessions, err := discovery.Discover(ctx)
                if err != nil {
                    return err
                }
                if rootJSON || flagJSON {
                    enc := json.NewEncoder(os.Stdout)
                    enc.SetIndent("", "  ")
                    return enc.Encode(sessions)
                }
                for _, s := range sessions {
                    fmt.Printf("%-8s %-20s %-30s pid=%d socket=%s\n", s.Status, s.Suffix, s.Prefix, s.PID, s.Socket)
                }
                return nil
            }
            // Default action: open the TUI
            return tui.Run()
        },
    }
    rootCmd.Flags().BoolVarP(&rootList, "list", "l", false, "list sessions and exit")
    rootCmd.Flags().BoolVar(&rootJSON, "json", false, "output JSON (with --list)")

    // list command
    listCmd := &cobra.Command{
        Use:   "list",
        Short: "List sessions",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            sessions, err := discovery.Discover(ctx)
            if err != nil {
                return err
            }
            if flagJSON {
                enc := json.NewEncoder(os.Stdout)
                enc.SetIndent("", "  ")
                return enc.Encode(sessions)
            }
            for _, s := range sessions {
                fmt.Printf("%-8s %-20s %-30s pid=%d socket=%s\n", s.Status, s.Suffix, s.Prefix, s.PID, s.Socket)
            }
            return nil
        },
    }
    listCmd.Flags().BoolVar(&flagJSON, "json", false, "output as JSON")

    // attach command (stub)
    attachCmd := &cobra.Command{
        Use:   "attach <selector>",
        Short: "Attach by suffix or folder prefix",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            selector := args[0]
            return discovery.AttachBySelector(selector)
        },
    }

    // start command (stub)
    var startName string
    var startDir string
    startCmd := &cobra.Command{
        Use:   "start [-- command...]",
        Short: "Create a new dtach session",
        RunE: func(cmd *cobra.Command, args []string) error {
            return discovery.StartSession(startDir, startName, args)
        },
    }
    startCmd.Flags().StringVarP(&startName, "name", "n", "", "suffix/name for the session")
    startCmd.Flags().StringVarP(&startDir, "chdir", "C", "", "launch directory")

    // rename command (stub)
    renameCmd := &cobra.Command{
        Use:   "rename <id|suffix> <new-suffix>",
        Short: "Rename session metadata (suffix)",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            return discovery.Rename(args[0], args[1])
        },
    }

    // kill command (stub)
    var killYes bool
    killCmd := &cobra.Command{
        Use:   "kill <id|suffix>",
        Short: "Kill a session's process",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            return discovery.Kill(args[0], killYes)
        },
    }
    killCmd.Flags().BoolVar(&killYes, "yes", false, "do not prompt for confirmation")

    // prune command (stub)
    var olderThan int
    var pruneAll bool
    pruneCmd := &cobra.Command{
        Use:   "prune",
        Short: "Remove stale sessions and metadata",
        RunE: func(cmd *cobra.Command, args []string) error {
            return discovery.Prune(pruneAll, olderThan)
        },
    }
    pruneCmd.Flags().BoolVar(&pruneAll, "all", false, "remove all stale entries")
    pruneCmd.Flags().IntVar(&olderThan, "older-than", 0, "remove stale entries older than N days")

    rootCmd.AddCommand(listCmd, attachCmd, startCmd, renameCmd, killCmd, pruneCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
