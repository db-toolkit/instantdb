package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
)

const metadataDir = ".instant-db"

// getMetadataDir returns the metadata directory path
func getMetadataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	
	dir := filepath.Join(home, metadataDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create metadata directory: %w", err)
	}
	
	return dir, nil
}

// SaveInstance saves instance metadata to disk
func SaveInstance(instance *types.Instance) error {
	dir, err := getMetadataDir()
	if err != nil {
		return err
	}
	
	path := filepath.Join(dir, fmt.Sprintf("%s.json", instance.ID))
	
	data, err := json.MarshalIndent(instance, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal instance: %w", err)
	}
	
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write instance file: %w", err)
	}
	
	return nil
}

// LoadInstance loads instance metadata from disk
func LoadInstance(instanceID string) (*types.Instance, error) {
	dir, err := getMetadataDir()
	if err != nil {
		return nil, err
	}
	
	path := filepath.Join(dir, fmt.Sprintf("%s.json", instanceID))
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read instance file: %w", err)
	}
	
	var instance types.Instance
	if err := json.Unmarshal(data, &instance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal instance: %w", err)
	}
	
	return &instance, nil
}

// RemoveInstance removes instance metadata from disk
func RemoveInstance(instanceID string) error {
	dir, err := getMetadataDir()
	if err != nil {
		return err
	}
	
	path := filepath.Join(dir, fmt.Sprintf("%s.json", instanceID))
	
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove instance file: %w", err)
	}
	
	return nil
}

// ListInstances returns all saved instances
func ListInstances() ([]*types.Instance, error) {
	dir, err := getMetadataDir()
	if err != nil {
		return nil, err
	}
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata directory: %w", err)
	}
	
	var instances []*types.Instance
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		
		var instance types.Instance
		if err := json.Unmarshal(data, &instance); err != nil {
			continue
		}
		
		instances = append(instances, &instance)
	}
	
	return instances, nil
}
