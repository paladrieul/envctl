package crypto

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	keystoreEnvVar  = "ENVCTL_PASSPHRASE"
	keystoreFileEnv = "ENVCTL_PASSPHRASE_FILE"
)

// ErrNoPassphrase is returned when no passphrase can be resolved.
var ErrNoPassphrase = errors.New("no passphrase found: set ENVCTL_PASSPHRASE or ENVCTL_PASSPHRASE_FILE")

// ResolvePassphrase attempts to load the encryption passphrase from
// environment variables or a passphrase file, in that order.
func ResolvePassphrase() (string, error) {
	if pass := os.Getenv(keystoreEnvVar); pass != "" {
		return pass, nil
	}

	if filePath := os.Getenv(keystoreFileEnv); filePath != "" {
		return readPassphraseFile(filePath)
	}

	defaultPath := filepath.Join(defaultConfigDir(), ".envctl_passphrase")
	if _, err := os.Stat(defaultPath); err == nil {
		return readPassphraseFile(defaultPath)
	}

	return "", ErrNoPassphrase
}

func readPassphraseFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func defaultConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".config", "envctl")
}
