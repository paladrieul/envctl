package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of a target's environment variables.
type Snapshot struct {
	Target    string            `json:"target"`
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label,omitempty"`
	Vars      map[string]string `json:"vars"`
}

// SaveSnapshot writes a snapshot of the given target to the snapshots directory.
func (s *Store) SaveSnapshot(target, label string) (*Snapshot, error) {
	vars, err := s.Load(target)
	if err != nil {
		return nil, fmt.Errorf("load target %q: %w", target, err)
	}
	if len(vars) == 0 {
		return nil, fmt.Errorf("target %q has no variables to snapshot", target)
	}

	snap := &Snapshot{
		Target:    target,
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Vars:      vars,
	}

	dir := filepath.Join(s.dir, "snapshots", target)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("create snapshot dir: %w", err)
	}

	filename := snap.CreatedAt.Format("20060102T150405Z") + ".json"
	path := filepath.Join(dir, filename)

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal snapshot: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return nil, fmt.Errorf("write snapshot: %w", err)
	}

	return snap, nil
}

// ListSnapshots returns all snapshots for the given target, ordered oldest first.
func (s *Store) ListSnapshots(target string) ([]*Snapshot, error) {
	dir := filepath.Join(s.dir, "snapshots", target)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read snapshot dir: %w", err)
	}

	var snapshots []*Snapshot
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("read snapshot %q: %w", e.Name(), err)
		}
		var snap Snapshot
		if err := json.Unmarshal(data, &snap); err != nil {
			return nil, fmt.Errorf("parse snapshot %q: %w", e.Name(), err)
		}
		snapshots = append(snapshots, &snap)
	}
	return snapshots, nil
}

// RestoreSnapshot overwrites the target's variables with those from the given snapshot.
func (s *Store) RestoreSnapshot(snap *Snapshot) error {
	return s.Save(snap.Target, snap.Vars)
}
