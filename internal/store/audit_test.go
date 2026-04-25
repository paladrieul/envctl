package store

import (
	"os"
	"testing"
	"time"
)

func TestAuditAppendAndReadAll(t *testing.T) {
	dir := t.TempDir()
	al := NewAuditLog(dir)

	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		Target:    "production",
		Action:    ActionSet,
		Key:       "DB_PASSWORD",
		Encrypted: true,
	}
	if err := al.Append(entry); err != nil {
		t.Fatalf("Append: %v", err)
	}

	entries, err := al.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "DB_PASSWORD" {
		t.Errorf("expected key DB_PASSWORD, got %s", entries[0].Key)
	}
	if entries[0].Action != ActionSet {
		t.Errorf("expected action set, got %s", entries[0].Action)
	}
	if !entries[0].Encrypted {
		t.Error("expected encrypted=true")
	}
}

func TestAuditReadAllEmpty(t *testing.T) {
	dir := t.TempDir()
	al := NewAuditLog(dir)

	entries, err := al.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on missing file: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestAuditMultipleEntries(t *testing.T) {
	dir := t.TempDir()
	al := NewAuditLog(dir)

	actions := []AuditAction{ActionSet, ActionDelete, ActionImport}
	for _, action := range actions {
		if err := al.Append(AuditEntry{Target: "staging", Action: action, Key: "KEY"}); err != nil {
			t.Fatalf("Append %s: %v", action, err)
		}
	}

	entries, err := al.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i, e := range entries {
		if e.Action != actions[i] {
			t.Errorf("entry %d: expected action %s, got %s", i, actions[i], e.Action)
		}
		if e.Timestamp.IsZero() {
			t.Errorf("entry %d: timestamp was not set", i)
		}
	}
}

func TestAuditTimestampAutoSet(t *testing.T) {
	dir := t.TempDir()
	al := NewAuditLog(dir)

	before := time.Now().UTC()
	if err := al.Append(AuditEntry{Target: "dev", Action: ActionCopy, Key: "X"}); err != nil {
		t.Fatalf("Append: %v", err)
	}
	after := time.Now().UTC()

	entries, _ := al.ReadAll()
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}

func TestAuditUnreadableFile(t *testing.T) {
	dir := t.TempDir()
	al := NewAuditLog(dir)
	// Write a valid entry first, then corrupt permissions.
	_ = al.Append(AuditEntry{Target: "t", Action: ActionSet, Key: "K"})
	if err := os.Chmod(al.path, 0000); err != nil {
		t.Skip("cannot chmod, skipping")
	}
	t.Cleanup(func() { os.Chmod(al.path, 0600) })
	_, err := al.ReadAll()
	if err == nil {
		t.Error("expected error reading unreadable file")
	}
}
