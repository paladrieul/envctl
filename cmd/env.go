package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourorg/envctl/internal/crypto"
	"github.com/yourorg/envctl/internal/store"
)

var (
	envTarget  string
	storeDir   string
	decryptOut bool
)

func init() {
	envCmd.PersistentFlags().StringVarP(&envTarget, "target", "t", "default", "deployment target name")
	envCmd.PersistentFlags().StringVar(&storeDir, "store", ".envctl", "directory to store env files")

	envSetCmd.Flags().BoolVar(&decryptOut, "decrypt", false, "store value as plaintext (no encryption)")

	envCmd.AddCommand(envSetCmd, envGetCmd, envDelCmd, envListCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment variables for a target",
}

var envSetCmd = &cobra.Command{
	Use:   "set KEY=VALUE [KEY=VALUE ...]",
	Short: "Set one or more environment variables (encrypted by default)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s := store.New(storeDir)
		ef, err := s.Load(envTarget)
		if err != nil {
			return err
		}
		pass, err := crypto.ResolvePassphrase()
		if err != nil {
			return err
		}
		for _, arg := range args {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid format %q, expected KEY=VALUE", arg)
			}
			key, val := parts[0], parts[1]
			if !decryptOut {
				val, err = crypto.Encrypt(val, pass)
				if err != nil {
					return fmt.Errorf("encrypt %s: %w", key, err)
				}
				val = "enc:" + val
			}
			ef.Entries[key] = val
		}
		return s.Save(ef)
	},
}

var envGetCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Get and decrypt an environment variable",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s := store.New(storeDir)
		ef, err := s.Load(envTarget)
		if err != nil {
			return err
		}
		val, ok := ef.Entries[args[0]]
		if !ok {
			return fmt.Errorf("key %q not found in target %q", args[0], envTarget)
		}
		if strings.HasPrefix(val, "enc:") {
			pass, err := crypto.ResolvePassphrase()
			if err != nil {
				return err
			}
			val, err = crypto.Decrypt(strings.TrimPrefix(val, "enc:"), pass)
			if err != nil {
				return fmt.Errorf("decrypt %s: %w", args[0], err)
			}
		}
		fmt.Fprintln(os.Stdout, val)
		return nil
	},
}

var envDelCmd = &cobra.Command{
	Use:   "del KEY",
	Short: "Delete an environment variable from a target",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return store.New(storeDir).Delete(envTarget, args[0])
	},
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all keys for a target (values masked)",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := store.New(storeDir)
		ef, err := s.Load(envTarget)
		if err != nil {
			return err
		}
		keys := make([]string, 0, len(ef.Entries))
		for k := range ef.Entries {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := ef.Entries[k]
			if strings.HasPrefix(v, "enc:") {
				v = "[encrypted]"
			}
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	},
}
