package types

import "time"

type Status string

const (
    StatusUnknown  Status = "unknown"
    StatusDetached Status = "detached"
    StatusAttached Status = "attached"
    StatusStale    Status = "stale"
)

type Session struct {
    ID        string    `json:"id"`
    Socket    string    `json:"socket"`
    Prefix    string    `json:"prefix"`
    Suffix    string    `json:"suffix"`
    PID       int       `json:"pid"`
    Status    Status    `json:"status"`
    Command   string    `json:"command"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

