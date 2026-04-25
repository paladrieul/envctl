package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Tags maps environment variable keys to a list of string tags.
type Tags map[string][]string

func tagsPath(dir, target string) string {
	return filepath.Join(dir, target+".tags.json")
}

// LoadTags reads the tag metadata for a target. Returns empty Tags if none exist.
func (s *Store) LoadTags(target string) (Tags, error) {
	path := tagsPath(s.dir, target)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Tags{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read tags: %w", err)
	}
	var t Tags
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parse tags: %w", err)
	}
	return t, nil
}

// SaveTags persists the tag metadata for a target.
func (s *Store) SaveTags(target string, t Tags) error {
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal tags: %w", err)
	}
	return os.WriteFile(tagsPath(s.dir, target), data, 0o600)
}

// AddTag appends tag to the list for key in target (no duplicates).
func (s *Store) AddTag(target, key, tag string) error {
	t, err := s.LoadTags(target)
	if err != nil {
		return err
	}
	for _, existing := range t[key] {
		if existing == tag {
			return nil
		}
	}
	t[key] = append(t[key], tag)
	sort.Strings(t[key])
	return s.SaveTags(target, t)
}

// RemoveTag deletes tag from the list for key in target.
func (s *Store) RemoveTag(target, key, tag string) error {
	t, err := s.LoadTags(target)
	if err != nil {
		return err
	}
	filtered := t[key][:0]
	for _, existing := range t[key] {
		if existing != tag {
			filtered = append(filtered, existing)
		}
	}
	t[key] = filtered
	return s.SaveTags(target, t)
}
