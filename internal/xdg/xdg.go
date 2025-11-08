package xdg

import (
    "os"
    "path/filepath"
    "runtime"
)

func ConfigDir() string {
    if runtime.GOOS == "darwin" {
        home, _ := os.UserHomeDir()
        return filepath.Join(home, "Library", "Application Support", "tache")
    }
    if d, err := os.UserConfigDir(); err == nil {
        return filepath.Join(d, "tache")
    }
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".config", "tache")
}

func StateDir() string {
    if runtime.GOOS == "darwin" {
        home, _ := os.UserHomeDir()
        return filepath.Join(home, "Library", "Application Support", "tache", "state")
    }
    if d, err := os.UserHomeDir(); err == nil {
        return filepath.Join(d, ".local", "state", "tache")
    }
    return "."
}

