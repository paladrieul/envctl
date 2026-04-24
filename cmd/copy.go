package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envctl/internal/store"
)

func init() {
	var overwrite bool

	copyCmd := &cobra.Command{
		Use:   "copy <src-target> <dst-target>",
		Short: "Copy environment variables from one target to another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCopy(cmd, args[0], args[1], overwrite)
		},
	}

	copyCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false,
		"Overwrite existing keys in the destination target")

	NewRootCmd().AddCommand(copyCmd)
}

func runCopy(cmd *cobra.Command, src, dst string, overwrite bool) error {
	storeDir, err := cmd.Flags().GetString("store")
	if err != nil {
		return err
	}

	s, err := store.New(storeDir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	result, err := store.Copy(s, src, dst, overwrite)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(),
		"Copied %d key(s) from %q to %q", result.Copied, src, dst)
	if result.Skipped > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), " (%d skipped, use --overwrite to replace)", result.Skipped)
	}
	fmt.Fprintln(cmd.OutOrStdout())
	return nil
}
