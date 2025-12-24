// +build !windows

package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// IsProcessRunning checks if a process with the given PID is running
func IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// KillProcessOnPort finds and kills the process listening on a specific port
func KillProcessOnPort(port int) error {
	// Use lsof to find process on port
	cmd := exec.Command("lsof", "-ti", fmt.Sprintf(":%d", port))
	output, err := cmd.Output()
	if err != nil {
		// No process found on port
		return nil
	}

	pidStr := strings.TrimSpace(string(output))
	if pidStr == "" {
		return nil
	}

	// Split by newlines in case multiple processes
	pids := strings.Split(pidStr, "\n")
	
	for _, pidLine := range pids {
		pidLine = strings.TrimSpace(pidLine)
		if pidLine == "" {
			continue
		}

		pid, err := strconv.Atoi(pidLine)
		if err != nil {
			continue // Skip invalid PIDs
		}

		// Kill the process
		process, err := os.FindProcess(pid)
		if err != nil {
			continue
		}

		// Send SIGTERM for graceful shutdown
		process.Signal(syscall.SIGTERM)
	}

	return nil
}
