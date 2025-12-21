package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
	sqle "github.com/dolthub/go-mysql-server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: mysql-daemon <port>")
		os.Exit(1)
	}

	port := os.Args[1]
	
	db := memory.NewDatabase("mysql")
	pro := memory.NewDBProvider(db)
	engine := sqle.NewDefault(pro)
	engine.Analyzer.Catalog.InfoSchema = information_schema.NewInformationSchemaDatabase()
	
	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("127.0.0.1:%s", port),
	}
	
	s, err := server.NewServer(config, engine, nil, memory.NewSessionBuilder(pro), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		s.Close()
		cancel()
	}()

	if err := s.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}

	<-ctx.Done()
}
