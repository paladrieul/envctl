package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

// rootCmd is the base command for the envctl CLI.
var rootCmd = &cobra.Command{
	Use:   "envctl",
	Short: "Manage environment variables across deployment targets",
	Long: `envctl is a CLI tool for managing environment variables across
multiple deployment targets with encryption support.

Securely store, retrieve, and sync environment variables for
development, staging, and production environments.`,
	SilenceUsage: true,
}

// versionCmd prints the current version of envctl.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of envctl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("envctl version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Persistent flags available to all subcommands
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to config file (default: ~/.envctl/config.yaml)")
	rootCmd.PersistentFlags().StringP("target", "t", "", "deployment target (e.g. dev, staging, production)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
