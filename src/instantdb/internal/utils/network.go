package utils

import (
	"fmt"
	"net"
)

// GetFreePort finds and returns an available port
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("failed to resolve address: %w", err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("failed to listen: %w", err)
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port, nil
}
