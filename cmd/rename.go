package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envctl/internal/store"
)

func init() {
	var target string
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "rename <old-key> <new-key>",
		Short: "Rename an environment variable key within a target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRename(cmd, args, target, overwrite)
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "default", "deployment target")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite destination key if it already exists")

	NewRootCmd().AddCommand(cmd)
}

func runRename(cmd *cobra.Command, args []string, target string, overwrite bool) error {
	oldKey := args[0]
	newKey := args[1]

	s, err := store.New("")
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	res, err := store.Rename(s, target, oldKey, newKey, overwrite)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Renamed %q → %q in target %q\n", res.OldVal, res.Key, target)
	return nil
}
