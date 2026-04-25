package store

import (
	"testing"
)

func TestLoadTagsEmpty(t *testing.T) {
	s := tempStore(t)
	tags, err := s.LoadTags("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("expected empty tags, got %v", tags)
	}
}

func TestAddAndLoadTag(t *testing.T) {
	s := tempStore(t)
	if err := s.AddTag("prod", "DB_URL", "sensitive"); err != nil {
		t.Fatalf("AddTag: %v", err)
	}
	tags, err := s.LoadTags("prod")
	if err != nil {
		t.Fatalf("LoadTags: %v", err)
	}
	if len(tags["DB_URL"]) != 1 || tags["DB_URL"][0] != "sensitive" {
		t.Fatalf("expected [sensitive], got %v", tags["DB_URL"])
	}
}

func TestAddTagNoDuplicates(t *testing.T) {
	s := tempStore(t)
	_ = s.AddTag("prod", "API_KEY", "secret")
	_ = s.AddTag("prod", "API_KEY", "secret")
	tags, _ := s.LoadTags("prod")
	if len(tags["API_KEY"]) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags["API_KEY"]))
	}
}

func TestRemoveTag(t *testing.T) {
	s := tempStore(t)
	_ = s.AddTag("prod", "TOKEN", "sensitive")
	_ = s.AddTag("prod", "TOKEN", "rotate")
	if err := s.RemoveTag("prod", "TOKEN", "sensitive"); err != nil {
		t.Fatalf("RemoveTag: %v", err)
	}
	tags, _ := s.LoadTags("prod")
	if len(tags["TOKEN"]) != 1 || tags["TOKEN"][0] != "rotate" {
		t.Fatalf("expected [rotate], got %v", tags["TOKEN"])
	}
}

func TestRemoveTagNotPresent(t *testing.T) {
	s := tempStore(t)
	_ = s.AddTag("prod", "KEY", "info")
	if err := s.RemoveTag("prod", "KEY", "nonexistent"); err != nil {
		t.Fatalf("RemoveTag on missing tag should not error: %v", err)
	}
	tags, _ := s.LoadTags("prod")
	if len(tags["KEY"]) != 1 {
		t.Fatalf("tag list should be unchanged")
	}
}

func TestTagsPersistedAcrossLoad(t *testing.T) {
	s := tempStore(t)
	_ = s.AddTag("staging", "DB_PASS", "sensitive")
	_ = s.AddTag("staging", "DB_PASS", "managed")
	tags, err := s.LoadTags("staging")
	if err != nil {
		t.Fatalf("LoadTags: %v", err)
	}
	if len(tags["DB_PASS"]) != 2 {
		t.Fatalf("expected 2 tags, got %v", tags["DB_PASS"])
	}
}
