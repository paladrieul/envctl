package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envctl/internal/store"
)

func tempStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	return store.New(dir)
}

func TestLoadNonExistentReturnsEmpty(t *testing.T) {
	s := tempStore(t)
	ef, err := s.Load("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 0 {
		t.Errorf("expected empty entries, got %v", ef.Entries)
	}
	if ef.Target != "staging" {
		t.Errorf("expected target=staging, got %s", ef.Target)
	}
}

func TestSaveAndLoad(t *testing.T) {
	s := tempStore(t)
	ef := &store.EnvFile{
		Target:  "production",
		Version: 0,
		Entries: map[string]string{"DB_URL": "enc:abc123", "API_KEY": "enc:xyz"},
	}
	if err := s.Save(ef); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := s.Load("production")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Entries["DB_URL"] != "enc:abc123" {
		t.Errorf("unexpected DB_URL: %s", loaded.Entries["DB_URL"])
	}
	if loaded.Version != 1 {
		t.Errorf("expected version=1 after save, got %d", loaded.Version)
	}
}

func TestDeleteKey(t *testing.T) {
	s := tempStore(t)
	ef := &store.EnvFile{
		Target:  "dev",
		Entries: map[string]string{"FOO": "bar", "BAZ": "qux"},
	}
	_ = s.Save(ef)
	if err := s.Delete("dev", "FOO"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	loaded, _ := s.Load("dev")
	if _, ok := loaded.Entries["FOO"]; ok {
		t.Error("expected FOO to be deleted")
	}
	if loaded.Entries["BAZ"] != "qux" {
		t.Error("expected BAZ to remain")
	}
}

func TestListTargets(t *testing.T) {
	s := tempStore(t)
	for _, tgt := range []string{"dev", "staging", "prod"} {
		_ = s.Save(&store.EnvFile{Target: tgt, Entries: map[string]string{}})
	}
	// add a non-json file that should be ignored
	_ = os.WriteFile(filepath.Join(s.BaseDir, "notes.txt"), []byte("hi"), 0600)

	targets, err := s.ListTargets()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(targets) != 3 {
		t.Errorf("expected 3 targets, got %d: %v", len(targets), targets)
	}
}
