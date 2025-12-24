package test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestPostgresStart(t *testing.T) {
	ctx := context.Background()
	engine := setupTestEngine(t, "postgres")
	
	config := createTestConfig("test-postgres", false)
	instance, err := engine.Start(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start postgres: %v", err)
	}
	defer cleanupInstance(t, engine, instance.ID)
	
	if instance.Engine != "postgres" {
		t.Errorf("Expected engine postgres, got %s", instance.Engine)
	}
	if instance.Port == 0 {
		t.Error("Port should be assigned")
	}
	
	t.Logf("Started postgres on port %d", instance.Port)
}

func TestPostgresConnection(t *testing.T) {
	ctx := context.Background()
	engine := setupTestEngine(t, "postgres")
	
	config := createTestConfig("test-postgres-conn", false)
	instance, err := engine.Start(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start postgres: %v", err)
	}
	defer cleanupInstance(t, engine, instance.ID)
	
	connURL, err := engine.GetConnectionURL(instance.ID)
	if err != nil {
		t.Fatalf("Failed to get connection URL: %v", err)
	}
	
	time.Sleep(2 * time.Second)
	
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		t.Fatalf("Failed to open connection: %v", err)
	}
	defer db.Close()
	
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping postgres: %v", err)
	}
	
	t.Log("Successfully connected to postgres")
}

func TestPostgresPauseResume(t *testing.T) {
	ctx := context.Background()
	engine := setupTestEngine(t, "postgres")
	
	config := createTestConfig("test-postgres-pause", false)
	instance, err := engine.Start(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start postgres: %v", err)
	}
	defer cleanupInstance(t, engine, instance.ID)
	
	if err := engine.Pause(ctx, instance.ID); err != nil {
		t.Fatalf("Failed to pause: %v", err)
	}
	
	status, err := engine.Status(ctx, instance.ID)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}
	if status.Running {
		t.Error("Instance should not be running after pause")
	}
	
	if err := engine.Resume(ctx, instance.ID); err != nil {
		t.Fatalf("Failed to resume: %v", err)
	}
	
	time.Sleep(2 * time.Second)
	
	status, err = engine.Status(ctx, instance.ID)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}
	if !status.Running {
		t.Error("Instance should be running after resume")
	}
	
	t.Log("Pause/resume successful")
}

func TestPostgresPersistence(t *testing.T) {
	ctx := context.Background()
	engine := setupTestEngine(t, "postgres")
	
	config := createTestConfig("test-postgres-persist", true)
	instance, err := engine.Start(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start postgres: %v", err)
	}
	
	dataDir := instance.DataDir
	
	if err := engine.Stop(ctx, instance.ID); err != nil {
		t.Fatalf("Failed to stop: %v", err)
	}
	
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Error("Data directory should persist after stop")
	} else {
		os.RemoveAll(dataDir)
	}
	
	t.Log("Persistence test passed")
}
