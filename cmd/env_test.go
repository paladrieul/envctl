package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yourorg/envctl/cmd"
)

func runCmd(t *testing.T, root *cobra.Command, args ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestEnvSetAndGet(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVCTL_PASSPHRASE", "test-passphrase-123")

	root := cmd.NewRootCmd()
	_, err := runCmd(t, root, "env", "set", "--store", dir, "--target", "dev", "SECRET=mysecret")
	if err != nil {
		t.Fatalf("set: %v", err)
	}

	out, err := runCmd(t, root, "env", "get", "--store", dir, "--target", "dev", "SECRET")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if strings.TrimSpace(out) != "mysecret" {
		t.Errorf("expected mysecret, got %q", out)
	}
}

func TestEnvList(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVCTL_PASSPHRASE", "test-passphrase-123")

	root := cmd.NewRootCmd()
	_, _ = runCmd(t, root, "env", "set", "--store", dir, "--target", "staging", "FOO=bar", "BAZ=qux")

	out, err := runCmd(t, root, "env", "list", "--store", dir, "--target", "staging")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "BAZ=[encrypted]") {
		t.Errorf("expected encrypted marker in list, got:\n%s", out)
	}
	if !strings.Contains(out, "FOO=[encrypted]") {
		t.Errorf("expected FOO in list, got:\n%s", out)
	}
}

func TestEnvDel(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVCTL_PASSPHRASE", "test-passphrase-123")

	root := cmd.NewRootCmd()
	_, _ = runCmd(t, root, "env", "set", "--store", dir, "--target", "prod", "REMOVE_ME=val")
	_, err := runCmd(t, root, "env", "del", "--store", dir, "--target", "prod", "REMOVE_ME")
	if err != nil {
		t.Fatalf("del: %v", err)
	}
	out, _ := runCmd(t, root, "env", "list", "--store", dir, "--target", "prod")
	if strings.Contains(out, "REMOVE_ME") {
		t.Errorf("expected REMOVE_ME to be gone, got:\n%s", out)
	}
}

func TestEnvSetPlaintext(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVCTL_PASSPHRASE", "test-passphrase-123")

	root := cmd.NewRootCmd()
	_, err := runCmd(t, root, "env", "set", "--store", dir, "--target", "dev", "--decrypt", "PLAIN=hello")
	if err != nil {
		t.Fatalf("set plaintext: %v", err)
	}
	out, _ := runCmd(t, root, "env", "list", "--store", dir, "--target", "dev")
	if !strings.Contains(out, "PLAIN=hello") {
		t.Errorf("expected PLAIN=hello in list, got:\n%s", out)
	}
}

func TestEnvGetMissingKey(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVCTL_PASSPHRASE", "test-passphrase-123")
	root := cmd.NewRootCmd()
	_, err := runCmd(t, root, "env", "get", "--store", dir, "--target", "dev", "MISSING")
	if err == nil {
		t.Error("expected error for missing key")
	}
	_ = os.Getenv("")
}
