package cmd

import (
	"strings"
	"testing"
)

// TestSnapshotMultipleSaves verifies that multiple snapshots accumulate correctly.
func TestSnapshotMultipleSaves(t *testing.T) {
	dir := t.TempDir()

	for _, kv := range []string{"A=1", "B=2", "C=3"} {
		if _, err := runCmd(dir, "env", "set", "prod", kv); err != nil {
			t.Fatalf("env set %s: %v", kv, err)
		}
	}

	for i, label := range []string{"snap1", "snap2", "snap3"} {
		if _, err := runCmd(dir, "snapshot", "save", "prod", "--label", label); err != nil {
			t.Fatalf("snapshot save %d: %v", i, err)
		}
	}

	out, err := runCmd(dir, "snapshot", "list", "prod")
	if err != nil {
		t.Fatalf("snapshot list: %v", err)
	}
	for _, label := range []string{"snap1", "snap2", "snap3"} {
		if !strings.Contains(out, label) {
			t.Errorf("expected label %q in list output:\n%s", label, out)
		}
	}
}

// TestSnapshotSaveEmptyTargetFails ensures saving a snapshot of a non-existent target fails.
func TestSnapshotSaveEmptyTargetFails(t *testing.T) {
	dir := t.TempDir()

	_, err := runCmd(dir, "snapshot", "save", "nonexistent")
	if err == nil {
		t.Error("expected error when snapshotting empty target, got nil")
	}
}

// TestSnapshotRestoreSecondSnapshot verifies restoring by index selects the correct snapshot.
func TestSnapshotRestoreSecondSnapshot(t *testing.T) {
	dir := t.TempDir()

	// First state
	if _, err := runCmd(dir, "env", "set", "qa", "VER=1.0"); err != nil {
		t.Fatalf("set VER 1.0: %v", err)
	}
	if _, err := runCmd(dir, "snapshot", "save", "qa", "--label", "first"); err != nil {
		t.Fatalf("save first: %v", err)
	}

	// Second state
	if _, err := runCmd(dir, "env", "set", "qa", "VER=2.0"); err != nil {
		t.Fatalf("set VER 2.0: %v", err)
	}
	if _, err := runCmd(dir, "snapshot", "save", "qa", "--label", "second"); err != nil {
		t.Fatalf("save second: %v", err)
	}

	// Advance to third state
	if _, err := runCmd(dir, "env", "set", "qa", "VER=3.0"); err != nil {
		t.Fatalf("set VER 3.0: %v", err)
	}

	// Restore from snapshot 2
	if _, err := runCmd(dir, "snapshot", "restore", "qa", "2"); err != nil {
		t.Fatalf("restore snapshot 2: %v", err)
	}

	out, err := runCmd(dir, "env", "get", "qa", "VER")
	if err != nil {
		t.Fatalf("env get VER: %v", err)
	}
	if !strings.Contains(out, "2.0") {
		t.Errorf("expected VER=2.0 after restoring snapshot 2, got: %s", out)
	}
}
