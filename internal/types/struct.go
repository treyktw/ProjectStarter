package types

import (
	"os/exec"
	"time"
)

type projectTemplate struct {
	name     string
	initFunc func(string, string) (*exec.Cmd, error)
}

type ProjectStats struct {
	LastModified time.Time
	TotalSize    int64
	FileCount    int
}
