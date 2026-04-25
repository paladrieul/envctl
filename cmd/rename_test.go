package cmd

import (
	"strings"
	"testing"
)

func TestRenameBasic(t *testing.T) {
	// Set a key then rename it
	_, err := runCmd(t, "env", "set", "--target", "staging", "OLD_NAME=myvalue")
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}

	out, err := runCmd(t, "rename", "--target", "staging", "OLD_NAME", "NEW_NAME")
	if err != nil {
		t.Fatalf("rename failed: %v", err)
	}
	if !strings.Contains(out, "NEW_NAME") {
		t.Errorf("expected output to mention NEW_NAME, got: %s", out)
	}

	// Confirm old key gone and new key present
	out, err = runCmd(t, "env", "list", "--target", "staging")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if strings.Contains(out, "OLD_NAME") {
		t.Error("old key OLD_NAME should not appear after rename")
	}
	if !strings.Contains(out, "NEW_NAME") {
		t.Error("new key NEW_NAME should appear after rename")
	}
}

func TestRenameMissingKeyError(t *testing.T) {
	_, err := runCmd(t, "rename", "--target", "staging", "DOES_NOT_EXIST", "SOMETHING")
	if err == nil {
		t.Fatal("expected error when renaming a non-existent key")
	}
}

func TestRenameSameKeyError(t *testing.T) {
	_, err := runCmd(t, "env", "set", "--target", "staging", "MYKEY=val")
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}

	_, err = runCmd(t, "rename", "--target", "staging", "MYKEY", "MYKEY")
	if err == nil {
		t.Fatal("expected error when old and new key names are identical")
	}
}

func TestRenameOverwriteFlag(t *testing.T) {
	_, _ = runCmd(t, "env", "set", "--target", "staging", "SRC=srcval")
	_, _ = runCmd(t, "env", "set", "--target", "staging", "DST=dstval")

	_, err := runCmd(t, "rename", "--target", "staging", "SRC", "DST")
	if err == nil {
		t.Fatal("expected error without --overwrite flag")
	}

	_, err = runCmd(t, "rename", "--target", "staging", "--overwrite", "SRC", "DST")
	if err != nil {
		t.Fatalf("expected success with --overwrite, got: %v", err)
	}
}
