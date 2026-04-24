package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd builds and returns the root cobra command.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envctl",
		Short: "envctl — manage encrypted environment variables across deployment targets",
		SilenceUsage: true,
		SilenceErrors: true,
	}
	root.AddCommand(envCmd)
	return root
}

// Execute is the entry point called from main.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
