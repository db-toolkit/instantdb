package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/engines"
	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
)

func setupTestEngine(t *testing.T, engineType string) engines.Engine {
	homeDir, _ := os.UserHomeDir()
	baseDir := filepath.Join(homeDir, ".instant-db-test", "data")
	
	var engine engines.Engine
	switch engineType {
	case "postgres":
		engine = engines.NewPostgresEngine(baseDir)
	case "redis":
		engine = engines.NewRedisEngine(baseDir)
	case "mysql":
		engine = engines.NewMySQLEngine(baseDir)
	}
	
	return engine
}

func cleanupInstance(t *testing.T, engine engines.Engine, instanceID string) {
	ctx := context.Background()
	if err := engine.Stop(ctx, instanceID); err != nil {
		t.Logf("Cleanup warning: %v", err)
	}
}

func createTestConfig(name string, persist bool) types.Config {
	return types.Config{
		Name:     name,
		Port:     0,
		Persist:  persist,
		Username: "",
		Password: "",
	}
}
