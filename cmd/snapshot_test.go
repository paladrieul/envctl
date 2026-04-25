package cmd

import (
	"strings"
	"testing"
)

func TestSnapshotSaveAndList(t *testing.T) {
	dir := t.TempDir()

	// Set a variable first
	out, err := runCmd(dir, "env", "set", "prod", "DB_URL=postgres://localhost")
	if err != nil {
		t.Fatalf("env set: %v\n%s", err, out)
	}

	// Save snapshot
	out, err = runCmd(dir, "snapshot", "save", "prod", "--label", "v1")
	if err != nil {
		t.Fatalf("snapshot save: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Snapshot saved") {
		t.Errorf("expected 'Snapshot saved', got: %s", out)
	}

	// List snapshots
	out, err = runCmd(dir, "snapshot", "list", "prod")
	if err != nil {
		t.Fatalf("snapshot list: %v\n%s", err, out)
	}
	if !strings.Contains(out, "v1") {
		t.Errorf("expected label v1 in list output, got: %s", out)
	}
	if !strings.Contains(out, "1") {
		t.Errorf("expected var count in list output, got: %s", out)
	}
}

func TestSnapshotListEmpty(t *testing.T) {
	dir := t.TempDir()

	out, err := runCmd(dir, "snapshot", "list", "ghost")
	if err != nil {
		t.Fatalf("snapshot list: %v\n%s", err, out)
	}
	if !strings.Contains(out, "No snapshots found") {
		t.Errorf("expected no-snapshots message, got: %s", out)
	}
}

func TestSnapshotRestore(t *testing.T) {
	dir := t.TempDir()

	_, err := runCmd(dir, "env", "set", "staging", "API_KEY=original")
	if err != nil {
		t.Fatalf("env set original: %v", err)
	}

	_, err = runCmd(dir, "snapshot", "save", "staging")
	if err != nil {
		t.Fatalf("snapshot save: %v", err)
	}

	_, err = runCmd(dir, "env", "set", "staging", "API_KEY=changed")
	if err != nil {
		t.Fatalf("env set changed: %v", err)
	}

	out, err := runCmd(dir, "snapshot", "restore", "staging", "1")
	if err != nil {
		t.Fatalf("snapshot restore: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Restored") {
		t.Errorf("expected Restored in output, got: %s", out)
	}

	out, err = runCmd(dir, "env", "get", "staging", "API_KEY")
	if err != nil {
		t.Fatalf("env get after restore: %v", err)
	}
	if !strings.Contains(out, "original") {
		t.Errorf("expected original value after restore, got: %s", out)
	}
}

func TestSnapshotRestoreInvalidIndex(t *testing.T) {
	dir := t.TempDir()

	_, err := runCmd(dir, "env", "set", "dev", "X=1")
	if err != nil {
		t.Fatalf("env set: %v", err)
	}
	_, err = runCmd(dir, "snapshot", "save", "dev")
	if err != nil {
		t.Fatalf("snapshot save: %v", err)
	}

	_, err = runCmd(dir, "snapshot", "restore", "dev", "99")
	if err == nil {
		t.Error("expected error for out-of-range index, got nil")
	}
}
