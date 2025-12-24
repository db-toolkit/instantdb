package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// KillProcessOnPort finds and kills the process listening on a specific port (Windows)
func KillProcessOnPort(port int) error {
	// Use netstat to find process on port
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf(":%d", port)) && strings.Contains(line, "LISTENING") {
			fields := strings.Fields(line)
			if len(fields) > 4 {
				pid, err := strconv.Atoi(fields[len(fields)-1])
				if err != nil {
					continue
				}
				
				// Kill the process
				killCmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
				killCmd.Run()
			}
		}
	}

	return nil
}
