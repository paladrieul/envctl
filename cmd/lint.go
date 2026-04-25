package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/envctl/internal/store"
)

func init() {
	var storeDir string
	var target string
	var strict bool

	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Check environment variables for common issues",
		Long: `Lint inspects keys in a target for naming convention violations,
empty values, and case-insensitive duplicate keys.

Use --strict to exit with a non-zero status when any issues are found.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLint(storeDir, target, strict)
		},
	}

	cmd.Flags().StringVar(&storeDir, "store", "", "path to store directory (default: ~/.config/envctl)")
	cmd.Flags().StringVarP(&target, "target", "t", "default", "target environment to lint")
	cmd.Flags().BoolVar(&strict, "strict", false, "exit non-zero if any issues are found")

	NewRootCmd().AddCommand(cmd)
}

func runLint(storeDir, target string, strict bool) error {
	s, err := store.New(storeDir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	results, err := store.Lint(s, target)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Printf("✓ No issues found in target %q\n", target)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TARGET\tKEY\tISSUE")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\n", r.Target, r.Key, r.Issue)
	}
	w.Flush()

	if strict {
		return fmt.Errorf("lint: %d issue(s) found in target %q", len(results), target)
	}
	return nil
}
