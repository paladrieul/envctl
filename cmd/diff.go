package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envctl/internal/store"
)

func init() {
	diffCmd := &cobra.Command{
		Use:   "diff <target-a> <target-b>",
		Short: "Show differences between two targets",
		Args:  cobra.ExactArgs(2),
		RunE:  runDiff,
	}
	NewRootCmd().AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	targetA, targetB := args[0], args[1]

	s, err := store.New("")
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}

	result, err := s.Diff(targetA, targetB)
	if err != nil {
		return fmt.Errorf("computing diff: %w", err)
	}

	w := cmd.OutOrStdout()
	hasOutput := false

	for _, k := range result.OnlyInA {
		fmt.Fprintf(w, "- [%s] %s\n", targetA, k)
		hasOutput = true
	}

	for _, k := range result.OnlyInB {
		fmt.Fprintf(w, "+ [%s] %s\n", targetB, k)
		hasOutput = true
	}

	for k, vals := range result.Changed {
		fmt.Fprintf(w, "~ %s\n  %s: %s\n  %s: %s\n", k, targetA, vals[0], targetB, vals[1])
		hasOutput = true
	}

	if !hasOutput {
		fmt.Fprintf(w, "No differences between %s and %s.\n", targetA, targetB)
	}

	if len(result.OnlyInA) > 0 || len(result.OnlyInB) > 0 || len(result.Changed) > 0 {
		os.Exit(1)
	}
	return nil
}
