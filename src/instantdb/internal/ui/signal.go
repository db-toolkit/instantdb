package ui

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// SetupSignalHandler sets up graceful handling for Ctrl+C and other signals
func SetupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n")
		fmt.Println(MutedStyle.Render("👋 Goodbye!"))
		os.Exit(0)
	}()
}
