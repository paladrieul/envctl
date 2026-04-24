package store

import (
	"testing"
)

func TestCopyBasic(t *testing.T) {
	s := tempStore(t)

	_ = s.Save("staging", map[string]string{"FOO": "bar", "BAZ": "qux"})

	res, err := Copy(s, "staging", "production", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Copied != 2 || res.Skipped != 0 {
		t.Fatalf("expected 2 copied, 0 skipped; got %+v", res)
	}

	dst, _ := s.Load("production")
	if dst["FOO"] != "bar" || dst["BAZ"] != "qux" {
		t.Errorf("destination does not contain expected values: %v", dst)
	}
}

func TestCopyNoOverwrite(t *testing.T) {
	s := tempStore(t)

	_ = s.Save("staging", map[string]string{"FOO": "new", "EXTRA": "val"})
	_ = s.Save("production", map[string]string{"FOO": "existing"})

	res, err := Copy(s, "staging", "production", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Copied != 1 || res.Skipped != 1 {
		t.Fatalf("expected 1 copied, 1 skipped; got %+v", res)
	}

	dst, _ := s.Load("production")
	if dst["FOO"] != "existing" {
		t.Errorf("FOO should not have been overwritten, got %q", dst["FOO"])
	}
	if dst["EXTRA"] != "val" {
		t.Errorf("EXTRA should have been copied, got %q", dst["EXTRA"])
	}
}

func TestCopyWithOverwrite(t *testing.T) {
	s := tempStore(t)

	_ = s.Save("staging", map[string]string{"FOO": "new"})
	_ = s.Save("production", map[string]string{"FOO": "existing"})

	res, err := Copy(s, "staging", "production", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Copied != 1 || res.Skipped != 0 {
		t.Fatalf("expected 1 copied, 0 skipped; got %+v", res)
	}

	dst, _ := s.Load("production")
	if dst["FOO"] != "new" {
		t.Errorf("FOO should have been overwritten, got %q", dst["FOO"])
	}
}

func TestCopySameTargetErrors(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("staging", map[string]string{"K": "v"})

	_, err := Copy(s, "staging", "staging", false)
	if err == nil {
		t.Fatal("expected error when src == dst")
	}
}

func TestCopyEmptySourceErrors(t *testing.T) {
	s := tempStore(t)

	_, err := Copy(s, "nonexistent", "production", false)
	if err == nil {
		t.Fatal("expected error for missing source target")
	}
}
