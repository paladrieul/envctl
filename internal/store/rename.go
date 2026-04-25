package store

import "fmt"

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	Key    string
	OldVal string
	NewVal string
}

// Rename renames a key within a target, preserving its value.
// Returns an error if the source key does not exist or the destination key already exists
// (unless overwrite is true).
func Rename(s *Store, target, oldKey, newKey string, overwrite bool) (*RenameResult, error) {
	if oldKey == newKey {
		return nil, fmt.Errorf("old and new key names are identical: %q", oldKey)
	}

	envs, err := s.Load(target)
	if err != nil {
		return nil, fmt.Errorf("load target %q: %w", target, err)
	}

	val, ok := envs[oldKey]
	if !ok {
		return nil, fmt.Errorf("key %q not found in target %q", oldKey, target)
	}

	if _, exists := envs[newKey]; exists && !overwrite {
		return nil, fmt.Errorf("key %q already exists in target %q; use --overwrite to replace it", newKey, target)
	}

	oldVal := val
	delete(envs, oldKey)
	envs[newKey] = val

	if err := s.Save(target, envs); err != nil {
		return nil, fmt.Errorf("save target %q: %w", target, err)
	}

	return &RenameResult{Key: newKey, OldVal: oldKey, NewVal: oldVal}, nil
}
