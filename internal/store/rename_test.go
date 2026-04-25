package store

import (
	"testing"
)

func TestRenameBasic(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("prod", map[string]string{"OLD_KEY": "hello"})

	res, err := Rename(s, "prod", "OLD_KEY", "NEW_KEY", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Key != "NEW_KEY" {
		t.Errorf("expected new key NEW_KEY, got %s", res.Key)
	}

	envs, _ := s.Load("prod")
	if _, ok := envs["OLD_KEY"]; ok {
		t.Error("old key should have been removed")
	}
	if envs["NEW_KEY"] != "hello" {
		t.Errorf("expected value 'hello', got %q", envs["NEW_KEY"])
	}
}

func TestRenameSameKeyErrors(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("prod", map[string]string{"KEY": "val"})

	_, err := Rename(s, "prod", "KEY", "KEY", false)
	if err == nil {
		t.Fatal("expected error for identical key names")
	}
}

func TestRenameKeyNotFound(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("prod", map[string]string{"A": "1"})

	_, err := Rename(s, "prod", "MISSING", "NEW", false)
	if err == nil {
		t.Fatal("expected error for missing source key")
	}
}

func TestRenameNoOverwrite(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("prod", map[string]string{"A": "1", "B": "2"})

	_, err := Rename(s, "prod", "A", "B", false)
	if err == nil {
		t.Fatal("expected error when destination key exists and overwrite=false")
	}
}

func TestRenameWithOverwrite(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("prod", map[string]string{"A": "new_val", "B": "old_val"})

	_, err := Rename(s, "prod", "A", "B", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	envs, _ := s.Load("prod")
	if envs["B"] != "new_val" {
		t.Errorf("expected B=new_val after overwrite, got %q", envs["B"])
	}
	if _, ok := envs["A"]; ok {
		t.Error("old key A should have been removed")
	}
}
