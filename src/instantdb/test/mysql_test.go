package test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestMySQLStart(t *testing.T) {
	ctx := context.Background()
	engine := setupTestEngine(t, "mysql")
	
	config := createTestConfig("test-mysql", false)
	instance, err := engine.Start(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start mysql: %v", err)
	}
	defer cleanupInstance(t, engine, instance.ID)
	
	if instance.Engine != "mysql" {
		t.Errorf("Expected engine mysql, got %s", instance.Engine)
	}
	if instance.Port == 0 {
		t.Error("Port should be assigned")
	}
	
	t.Logf("Started mysql on port %d", instance.Port)
}

func TestMySQLConnection(t *testing.T) {
	ctx := context.Background()
	engine := setupTestEngine(t, "mysql")
	
	config := createTestConfig("test-mysql-conn", false)
	instance, err := engine.Start(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start mysql: %v", err)
	}
	defer cleanupInstance(t, engine, instance.ID)
	
	time.Sleep(3 * time.Second)
	
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/mysql", instance.Username, instance.Password, instance.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to open connection: %v", err)
	}
	defer db.Close()
	
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping mysql: %v", err)
	}
	
	var version string
	if err := db.QueryRow("SELECT VERSION()").Scan(&version); err != nil {
		t.Fatalf("Failed to query version: %v", err)
	}
	
	t.Logf("Successfully connected to mysql version: %s", version)
}

func TestMySQLPauseResume(t *testing.T) {
	ctx := context.Background()
	engine := setupTestEngine(t, "mysql")
	
	config := createTestConfig("test-mysql-pause", false)
	instance, err := engine.Start(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start mysql: %v", err)
	}
	defer cleanupInstance(t, engine, instance.ID)
	
	time.Sleep(3 * time.Second)
	
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
	
	time.Sleep(3 * time.Second)
	
	status, err = engine.Status(ctx, instance.ID)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}
	if !status.Running {
		t.Error("Instance should be running after resume")
	}
	
	t.Log("Pause/resume successful")
}

func TestMySQLPersistence(t *testing.T) {
	ctx := context.Background()
	engine := setupTestEngine(t, "mysql")
	
	config := createTestConfig("test-mysql-persist", true)
	instance, err := engine.Start(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start mysql: %v", err)
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
