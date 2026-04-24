package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envctl/internal/store"
)

func init() {
	var (
		target    string
		format    string
		overwrite bool
		filePath  string
	)

	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import environment variables from a file",
		Long: `Import key-value pairs into a target from a dotenv or JSON file.

Examples:
  envctl import --target prod --file .env
  envctl import --target staging --file vars.json --format json --overwrite`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if target == "" {
				return fmt.Errorf("--target is required")
			}

			var r *os.File
			var err error
			if filePath == "-" || filePath == "" {
				r = os.Stdin
			} else {
				r, err = os.Open(filePath)
				if err != nil {
					return fmt.Errorf("open file: %w", err)
				}
				defer r.Close()
			}

			if format == "" {
				format = "dotenv"
			}

			s, err := store.New("")
			if err != nil {
				return err
			}

			res, err := store.Import(s, target, r, store.ImportFormat(format), overwrite)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(),
				"Import complete: %d imported, %d skipped (target: %s)\n",
				res.Imported, res.Skipped, target)
			return nil
		},
	}

	importCmd.Flags().StringVarP(&target, "target", "t", "", "deployment target name (required)")
	importCmd.Flags().StringVarP(&filePath, "file", "f", "", "path to import file (default: stdin)")
	importCmd.Flags().StringVar(&format, "format", "dotenv", "import format: dotenv or json")
	importCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys")

	RootCmd.AddCommand(importCmd)
}
