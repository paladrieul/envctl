package store

import (
	"strings"
	"testing"
)

func TestLintCleanTarget(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("prod", map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "secret",
	})

	results, err := Lint(s, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no lint issues, got %d: %+v", len(results), results)
	}
}

func TestLintDetectsLowercaseKey(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("dev", map[string]string{
		"db_host": "localhost",
	})

	results, err := Lint(s, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsIssue(results, "db_host", "UPPER_SNAKE_CASE") {
		t.Errorf("expected UPPER_SNAKE_CASE issue for db_host, got: %+v", results)
	}
}

func TestLintDetectsEmptyValue(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("staging", map[string]string{
		"EMPTY_KEY": "",
	})

	results, err := Lint(s, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsIssue(results, "EMPTY_KEY", "empty") {
		t.Errorf("expected empty value issue, got: %+v", results)
	}
}

func TestLintDetectsLeadingUnderscore(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("dev", map[string]string{
		"_PRIVATE": "value",
	})

	results, err := Lint(s, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsIssue(results, "_PRIVATE", "leading or trailing underscore") {
		t.Errorf("expected underscore issue, got: %+v", results)
	}
}

func TestLintNonExistentTargetReturnsEmpty(t *testing.T) {
	s := tempStore(t)
	results, err := Lint(s, "ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for non-existent target, got %d", len(results))
	}
}

// containsIssue reports whether results contains an entry for key whose Issue
// field contains substr. If substr is empty, any entry matching key returns true.
func containsIssue(results []LintResult, key, substr string) bool {
	for _, r := range results {
		if r.Key == key {
			if substr == "" {
				return true
			}
			if strings.Contains(r.Issue, substr) {
				return true
			}
		}
	}
	return false
}
