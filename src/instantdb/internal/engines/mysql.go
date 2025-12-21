package engines

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
	"github.com/db-toolkit/instant-db/src/instantdb/internal/utils"
	sqle "github.com/dolthub/go-mysql-server"
)

type MySQLEngine struct {
	baseDir   string
	instances map[string]*server.Server
}

func NewMySQLEngine(baseDir string) *MySQLEngine {
	return &MySQLEngine{
		baseDir:   baseDir,
		instances: make(map[string]*server.Server),
	}
}

func (e *MySQLEngine) Start(ctx context.Context, config types.Config) (*types.Instance, error) {
	instanceID := utils.GenerateID()
	
	if config.Name == "" {
		config.Name = fmt.Sprintf("mysql-%s", instanceID[:8])
	}
	if config.Port == 0 {
		port, err := utils.GetFreePort()
		if err != nil {
			return nil, fmt.Errorf("failed to allocate port: %w", err)
		}
		config.Port = port
	}
	if config.DataDir == "" {
		config.DataDir = filepath.Join(e.baseDir, instanceID)
	}
	if config.Username == "" {
		config.Username = "root"
	}
	if config.Password == "" {
		config.Password = "password"
	}

	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db := memory.NewDatabase("mysql")
	pro := memory.NewDBProvider(db)
	engine := sqle.NewDefault(pro)
	engine.Analyzer.Catalog.InfoSchema = information_schema.NewInformationSchemaDatabase()
	
	serverConfig := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("127.0.0.1:%d", config.Port),
	}
	
	s, err := server.NewServer(serverConfig, engine, nil, memory.NewSessionBuilder(pro), nil)
	if err != nil {
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	go func() {
		s.Start()
	}()

	time.Sleep(500 * time.Millisecond)

	e.instances[instanceID] = s

	instance := &types.Instance{
		ID:        instanceID,
		Name:      config.Name,
		Engine:    "mysql",
		Port:      config.Port,
		DataDir:   config.DataDir,
		PID:       0,
		Status:    "running",
		CreatedAt: time.Now().Unix(),
		Persist:   config.Persist,
		Username:  config.Username,
		Password:  config.Password,
	}

	if err := utils.SaveInstance(instance); err != nil {
		s.Close()
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to save instance: %w", err)
	}

	return instance, nil
}

func (e *MySQLEngine) Stop(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	if s, ok := e.instances[instanceID]; ok {
		s.Close()
		delete(e.instances, instanceID)
	}

	if !instance.Persist {
		os.RemoveAll(instance.DataDir)
	}

	metaFile := filepath.Join(os.Getenv("HOME"), ".instant-db", instanceID+".json")
	os.Remove(metaFile)
	return nil
}

func (e *MySQLEngine) Pause(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	if s, ok := e.instances[instanceID]; ok {
		s.Close()
		delete(e.instances, instanceID)
	}

	instance.Status = "paused"
	instance.Paused = true
	return utils.SaveInstance(instance)
}

func (e *MySQLEngine) Resume(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	engine := memory.NewDatabase("mysql")
	pro := memory.NewDBProvider(engine)
	eng := sqle.NewDefault(pro)
	eng.Analyzer.Catalog.InfoSchema = information_schema.NewInformationSchemaDatabase()
	
	serverConfig := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("127.0.0.1:%d", instance.Port),
	}
	
	s, err := server.NewServer(serverConfig, eng, nil, memory.NewSessionBuilder(pro), nil)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	go func() {
		s.Start()
	}()

	time.Sleep(500 * time.Millisecond)

	e.instances[instanceID] = s

	instance.Status = "running"
	instance.Paused = false
	return utils.SaveInstance(instance)
}

func (e *MySQLEngine) GetStatus(ctx context.Context, instanceID string) (string, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return "", fmt.Errorf("instance not found: %w", err)
	}
	return instance.Status, nil
}

func (e *MySQLEngine) GetConnectionURL(instanceID string) (string, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return "", fmt.Errorf("instance not found: %w", err)
	}
	return fmt.Sprintf("mysql://%s:%s@127.0.0.1:%d/mysql", instance.Username, instance.Password, instance.Port), nil
}

func (e *MySQLEngine) Status(ctx context.Context, instanceID string) (*types.Status, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return nil, fmt.Errorf("instance not found: %w", err)
	}
	return &types.Status{
		Running: instance.Status == "running",
		Message: instance.Status,
	}, nil
}

func (e *MySQLEngine) List() ([]*types.Instance, error) {
	all, err := utils.ListInstances()
	if err != nil {
		return nil, err
	}
	var mysql []*types.Instance
	for _, inst := range all {
		if inst.Engine == "mysql" {
			mysql = append(mysql, inst)
		}
	}
	return mysql, nil
}
