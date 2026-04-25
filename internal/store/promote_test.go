package store_test

import (
	"testing"

	"github.com/user/envctl/internal/store"
)

func TestPromoteBasic(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("staging", map[string]string{"FOO": "bar", "BAZ": "qux"})
	_ = s.Save("prod", map[string]string{})

	results, err := store.Promote(s, "staging", "prod", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	prodVars, _ := s.Load("prod")
	if prodVars["FOO"] != "bar" || prodVars["BAZ"] != "qux" {
		t.Errorf("prod vars not promoted correctly: %v", prodVars)
	}
}

func TestPromoteNoOverwrite(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("staging", map[string]string{"FOO": "new", "ONLY_STAGING": "yes"})
	_ = s.Save("prod", map[string]string{"FOO": "existing"})

	results, err := store.Promote(s, "staging", "prod", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only ONLY_STAGING should be promoted; FOO should be skipped
	if len(results) != 1 || results[0].Key != "ONLY_STAGING" {
		t.Errorf("expected only ONLY_STAGING promoted, got %v", results)
	}

	prodVars, _ := s.Load("prod")
	if prodVars["FOO"] != "existing" {
		t.Errorf("FOO should not be overwritten, got %q", prodVars["FOO"])
	}
}

func TestPromoteWithOverwrite(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("staging", map[string]string{"FOO": "new"})
	_ = s.Save("prod", map[string]string{"FOO": "old"})

	results, err := store.Promote(s, "staging", "prod", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Overwrite {
		t.Errorf("expected overwrite result, got %v", results)
	}

	prodVars, _ := s.Load("prod")
	if prodVars["FOO"] != "new" {
		t.Errorf("FOO should be overwritten to 'new', got %q", prodVars["FOO"])
	}
}

func TestPromoteSameTargetErrors(t *testing.T) {
	s := tempStore(t)
	_, err := store.Promote(s, "staging", "staging", false)
	if err == nil {
		t.Fatal("expected error for same source and destination")
	}
}

func TestPromoteEmptySourceErrors(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("staging", map[string]string{})
	_, err := store.Promote(s, "staging", "prod", false)
	if err == nil {
		t.Fatal("expected error for empty source target")
	}
}
