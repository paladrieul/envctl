package crypto_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envctl/internal/crypto"
)

func TestResolvePassphraseFromEnv(t *testing.T) {
	t.Setenv("ENVCTL_PASSPHRASE", "env-passphrase")

	pass, err := crypto.ResolvePassphrase()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pass != "env-passphrase" {
		t.Fatalf("expected 'env-passphrase', got %q", pass)
	}
}

func TestResolvePassphraseFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	passFile := filepath.Join(tmpDir, "passphrase")

	if err := os.WriteFile(passFile, []byte("file-passphrase\n"), 0600); err != nil {
		t.Fatalf("failed to write passphrase file: %v", err)
	}

	os.Unsetenv("ENVCTL_PASSPHRASE")
	t.Setenv("ENVCTL_PASSPHRASE_FILE", passFile)

	pass, err := crypto.ResolvePassphrase()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pass != "file-passphrase" {
		t.Fatalf("expected 'file-passphrase', got %q", pass)
	}
}

func TestResolvePassphraseNotFound(t *testing.T) {
	os.Unsetenv("ENVCTL_PASSPHRASE")
	os.Unsetenv("ENVCTL_PASSPHRASE_FILE")

	// Point home to a temp dir so the default path doesn't accidentally exist
	t.Setenv("HOME", t.TempDir())

	_, err := crypto.ResolvePassphrase()
	if err == nil {
		t.Fatal("expected ErrNoPassphrase, got nil")
	}
	if err != crypto.ErrNoPassphrase {
		t.Fatalf("expected ErrNoPassphrase, got: %v", err)
	}
}
