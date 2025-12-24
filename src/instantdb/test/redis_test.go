package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestRedisStart(t *testing.T) {
	ctx := context.Background()
	engine := setupTestEngine(t, "redis")
	
	config := createTestConfig("test-redis", false)
	instance, err := engine.Start(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start redis: %v", err)
	}
	defer cleanupInstance(t, engine, instance.ID)
	
	if instance.Engine != "redis" {
		t.Errorf("Expected engine redis, got %s", instance.Engine)
	}
	if instance.Port == 0 {
		t.Error("Port should be assigned")
	}
	
	t.Logf("Started redis on port %d", instance.Port)
}

func TestRedisConnection(t *testing.T) {
	ctx := context.Background()
	engine := setupTestEngine(t, "redis")
	
	config := createTestConfig("test-redis-conn", false)
	instance, err := engine.Start(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start redis: %v", err)
	}
	defer cleanupInstance(t, engine, instance.ID)
	
	time.Sleep(1 * time.Second)
	
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("127.0.0.1:%d", instance.Port),
		Password: instance.Password,
	})
	defer client.Close()
	
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to ping redis: %v", err)
	}
	
	if err := client.Set(ctx, "test-key", "test-value", 0).Err(); err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}
	
	val, err := client.Get(ctx, "test-key").Result()
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}
	if val != "test-value" {
		t.Errorf("Expected test-value, got %s", val)
	}
	
	t.Log("Successfully connected to redis and performed operations")
}

func TestRedisPauseResume(t *testing.T) {
	ctx := context.Background()
	engine := setupTestEngine(t, "redis")
	
	config := createTestConfig("test-redis-pause", false)
	instance, err := engine.Start(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start redis: %v", err)
	}
	defer cleanupInstance(t, engine, instance.ID)
	
	time.Sleep(1 * time.Second)
	
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
	
	time.Sleep(1 * time.Second)
	
	status, err = engine.Status(ctx, instance.ID)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}
	if !status.Running {
		t.Error("Instance should be running after resume")
	}
	
	t.Log("Pause/resume successful")
}

func TestRedisPersistence(t *testing.T) {
	ctx := context.Background()
	engine := setupTestEngine(t, "redis")
	
	config := createTestConfig("test-redis-persist", true)
	instance, err := engine.Start(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start redis: %v", err)
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
