package store

import (
	"testing"
	"time"
)

func TestSaveAndListSnapshots(t *testing.T) {
	s := tempStore(t)

	vars := map[string]string{"KEY1": "val1", "KEY2": "val2"}
	if err := s.Save("prod", vars); err != nil {
		t.Fatalf("save: %v", err)
	}

	snap, err := s.SaveSnapshot("prod", "before-deploy")
	if err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}
	if snap.Target != "prod" {
		t.Errorf("expected target prod, got %q", snap.Target)
	}
	if snap.Label != "before-deploy" {
		t.Errorf("expected label before-deploy, got %q", snap.Label)
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}

	list, err := s.ListSnapshots("prod")
	if err != nil {
		t.Fatalf("ListSnapshots: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(list))
	}
	if list[0].Vars["KEY1"] != "val1" {
		t.Errorf("snapshot vars mismatch")
	}
}

func TestListSnapshotsEmpty(t *testing.T) {
	s := tempStore(t)
	list, err := s.ListSnapshots("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list != nil {
		t.Errorf("expected nil for non-existent target, got %v", list)
	}
}

func TestSaveSnapshotEmptyTargetErrors(t *testing.T) {
	s := tempStore(t)
	_, err := s.SaveSnapshot("ghost", "")
	if err == nil {
		t.Error("expected error for empty target, got nil")
	}
}

func TestRestoreSnapshot(t *testing.T) {
	s := tempStore(t)

	original := map[string]string{"DB_URL": "postgres://old"}
	if err := s.Save("prod", original); err != nil {
		t.Fatalf("save original: %v", err)
	}

	snap, err := s.SaveSnapshot("prod", "")
	if err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	// Overwrite with new data
	if err := s.Save("prod", map[string]string{"DB_URL": "postgres://new"}); err != nil {
		t.Fatalf("save updated: %v", err)
	}

	if err := s.RestoreSnapshot(snap); err != nil {
		t.Fatalf("RestoreSnapshot: %v", err)
	}

	restored, err := s.Load("prod")
	if err != nil {
		t.Fatalf("load restored: %v", err)
	}
	if restored["DB_URL"] != "postgres://old" {
		t.Errorf("expected original value, got %q", restored["DB_URL"])
	}
}

func TestMultipleSnapshotsOrdered(t *testing.T) {
	s := tempStore(t)
	if err := s.Save("dev", map[string]string{"X": "1"}); err != nil {
		t.Fatalf("save: %v", err)
	}

	for i := 0; i < 3; i++ {
		time.Sleep(1100 * time.Millisecond)
		if _, err := s.SaveSnapshot("dev", ""); err != nil {
			t.Fatalf("snapshot %d: %v", i, err)
		}
	}

	list, err := s.ListSnapshots("dev")
	if err != nil {
		t.Fatalf("ListSnapshots: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 snapshots, got %d", len(list))
	}
	for i := 1; i < len(list); i++ {
		if list[i].CreatedAt.Before(list[i-1].CreatedAt) {
			t.Errorf("snapshots not in order at index %d", i)
		}
	}
}
