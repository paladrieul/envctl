package cmd

import (
	"strings"
	"testing"
)

func TestTagAddAndList(t *testing.T) {
	dir := t.TempDir()
	// First set a key so the target exists conceptually
	out, err := runCmd(t, dir, "env", "set", "DB_URL", "postgres://localhost", "-t", "prod")
	if err != nil {
		t.Fatalf("env set: %v\n%s", err, out)
	}

	out, err = runCmd(t, dir, "tag", "add", "DB_URL", "sensitive", "-t", "prod")
	if err != nil {
		t.Fatalf("tag add: %v\n%s", err, out)
	}
	if !strings.Contains(out, "tagged DB_URL") {
		t.Fatalf("unexpected output: %s", out)
	}

	out, err = runCmd(t, dir, "tag", "list", "-t", "prod")
	if err != nil {
		t.Fatalf("tag list: %v\n%s", err, out)
	}
	if !strings.Contains(out, "DB_URL") || !strings.Contains(out, "sensitive") {
		t.Fatalf("expected DB_URL sensitive in output, got: %s", out)
	}
}

func TestTagRemove(t *testing.T) {
	dir := t.TempDir()
	_, _ = runCmd(t, dir, "env", "set", "API_KEY", "abc123", "-t", "staging")
	_, _ = runCmd(t, dir, "tag", "add", "API_KEY", "secret", "-t", "staging")
	_, _ = runCmd(t, dir, "tag", "add", "API_KEY", "rotate", "-t", "staging")

	out, err := runCmd(t, dir, "tag", "remove", "API_KEY", "secret", "-t", "staging")
	if err != nil {
		t.Fatalf("tag remove: %v\n%s", err, out)
	}

	out, err = runCmd(t, dir, "tag", "list", "-t", "staging")
	if err != nil {
		t.Fatalf("tag list: %v\n%s", err, out)
	}
	if strings.Contains(out, "secret") {
		t.Fatalf("secret should have been removed, got: %s", out)
	}
	if !strings.Contains(out, "rotate") {
		t.Fatalf("rotate should still be present, got: %s", out)
	}
}

func TestTagListEmpty(t *testing.T) {
	dir := t.TempDir()
	out, err := runCmd(t, dir, "tag", "list", "-t", "empty")
	if err != nil {
		t.Fatalf("tag list: %v\n%s", err, out)
	}
	if !strings.Contains(out, "no tags found") {
		t.Fatalf("expected 'no tags found', got: %s", out)
	}
}
