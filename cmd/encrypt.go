package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envctl/internal/crypto"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt [file]",
	Short: "Encrypt an env file in place",
	Args:  cobra.ExactArgs(1),
	RunE:  runEncrypt,
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt [file]",
	Short: "Decrypt an env file in place",
	Args:  cobra.ExactArgs(1),
	RunE:  runDecrypt,
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	passphrase, err := crypto.ResolvePassphrase()
	if err != nil {
		return fmt.Errorf("passphrase error: %w", err)
	}

	plaintext, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	encrypted, err := crypto.Encrypt(plaintext, passphrase)
	if err != nil {
		return fmt.Errorf("encrypting: %w", err)
	}

	if err := os.WriteFile(filePath+".enc", []byte(encrypted), 0600); err != nil {
		return fmt.Errorf("writing encrypted file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Encrypted → %s.enc\n", filePath)
	return nil
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	passphrase, err := crypto.ResolvePassphrase()
	if err != nil {
		return fmt.Errorf("passphrase error: %w", err)
	}

	encoded, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	plaintext, err := crypto.Decrypt(string(encoded), passphrase)
	if err != nil {
		return fmt.Errorf("decrypting: %w", err)
	}

	outPath := filePath
	if len(filePath) > 4 && filePath[len(filePath)-4:] == ".enc" {
		outPath = filePath[:len(filePath)-4]
	}

	if err := os.WriteFile(outPath+".dec", plaintext, 0600); err != nil {
		return fmt.Errorf("writing decrypted file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Decrypted → %s.dec\n", outPath)
	return nil
}
