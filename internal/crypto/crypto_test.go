package crypto_test

import (
	"testing"

	"github.com/yourorg/envctl/internal/crypto"
)

func TestEncryptDecryptRoundtrip(t *testing.T) {
	plaintext := []byte("SECRET_KEY=supersecret123")
	passphrase := "my-strong-passphrase"

	encrypted, err := crypto.Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if encrypted == string(plaintext) {
		t.Fatal("encrypted text should differ from plaintext")
	}

	decrypted, err := crypto.Decrypt(encrypted, passphrase)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestDecryptWrongPassphrase(t *testing.T) {
	plaintext := []byte("DB_PASSWORD=secret")

	encrypted, err := crypto.Encrypt(plaintext, "correct-passphrase")
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = crypto.Decrypt(encrypted, "wrong-passphrase")
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
}

func TestDecryptInvalidBase64(t *testing.T) {
	_, err := crypto.Decrypt("not-valid-base64!!!", "passphrase")
	if err == nil {
		t.Fatal("expected error for invalid base64 input")
	}
}

func TestEncryptProducesUniqueOutputs(t *testing.T) {
	plaintext := []byte("SAME_VALUE=hello")
	passphrase := "passphrase"

	enc1, err := crypto.Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("first Encrypt failed: %v", err)
	}

	enc2, err := crypto.Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("second Encrypt failed: %v", err)
	}

	// Due to random nonce, outputs should differ
	if enc1 == enc2 {
		t.Fatal("expected different ciphertexts for same plaintext due to random nonce")
	}
}
