package ui

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	cancelFunc context.CancelFunc
)

// SetupSignalHandler sets up graceful handling for Ctrl+C and other signals
func SetupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancelFunc = cancel

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n")
		fmt.Println(WarningStyle.Render("⚠️  Interrupt received, cleaning up..."))
		
		// Cancel any ongoing operations
		if cancelFunc != nil {
			cancelFunc()
		}
		
		// Give a moment for cleanup
		// Then exit
		fmt.Println(MutedStyle.Render("👋 Goodbye!"))
		os.Exit(130) // Standard exit code for Ctrl+C
	}()

	return ctx
}
