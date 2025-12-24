// +build windows

package utils

import (
	"os"
)

// IsProcessRunning checks if a process with the given PID is running (Windows)
func IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Windows, FindProcess always succeeds, so we can't reliably check
	return process != nil
}
