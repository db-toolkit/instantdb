package utils

import (
	"fmt"
)

// ResolveInstance resolves an instance name or ID to an instance ID
func ResolveInstance(nameOrID string) (string, error) {
	instances, err := ListInstances()
	if err != nil {
		return "", err
	}

	// Check if it's a name
	for _, inst := range instances {
		if inst.Name == nameOrID {
			return inst.ID, nil
		}
	}

	// Check if it's a valid ID
	for _, inst := range instances {
		if inst.ID == nameOrID {
			return inst.ID, nil
		}
	}

	return "", fmt.Errorf("instance not found: %s", nameOrID)
}
