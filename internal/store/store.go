package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// EnvFile represents a collection of environment variables for a target.
type EnvFile struct {
	Target  string            `json:"target"`
	Version int               `json:"version"`
	Entries map[string]string `json:"entries"`
}

// Store manages env files on disk.
type Store struct {
	BaseDir string
}

// New creates a Store rooted at baseDir.
func New(baseDir string) *Store {
	return &Store{BaseDir: baseDir}
}

// path returns the file path for a given target.
func (s *Store) path(target string) string {
	return filepath.Join(s.BaseDir, target+".json")
}

// Load reads the EnvFile for the given target.
func (s *Store) Load(target string) (*EnvFile, error) {
	data, err := os.ReadFile(s.path(target))
	if err != nil {
		if os.IsNotExist(err) {
			return &EnvFile{Target: target, Version: 1, Entries: map[string]string{}}, nil
		}
		return nil, fmt.Errorf("store: read %s: %w", target, err)
	}
	var ef EnvFile
	if err := json.Unmarshal(data, &ef); err != nil {
		return nil, fmt.Errorf("store: unmarshal %s: %w", target, err)
	}
	return &ef, nil
}

// Save writes the EnvFile to disk, creating directories as needed.
func (s *Store) Save(ef *EnvFile) error {
	if err := os.MkdirAll(s.BaseDir, 0700); err != nil {
		return fmt.Errorf("store: mkdir: %w", err)
	}
	ef.Version++
	data, err := json.MarshalIndent(ef, "", "  ")
	if err != nil {
		return fmt.Errorf("store: marshal: %w", err)
	}
	if err := os.WriteFile(s.path(ef.Target), data, 0600); err != nil {
		return fmt.Errorf("store: write %s: %w", ef.Target, err)
	}
	return nil
}

// Delete removes a key from the target env file and persists the change.
func (s *Store) Delete(target, key string) error {
	ef, err := s.Load(target)
	if err != nil {
		return err
	}
	delete(ef.Entries, key)
	return s.Save(ef)
}

// ListTargets returns all target names found in BaseDir.
func (s *Store) ListTargets() ([]string, error) {
	entries, err := os.ReadDir(s.BaseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("store: readdir: %w", err)
	}
	var targets []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			targets = append(targets, e.Name()[:len(e.Name())-5])
		}
	}
	return targets, nil
}
