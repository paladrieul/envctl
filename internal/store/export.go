package store

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportFormat represents the output format for environment variable export.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
	FormatExport ExportFormat = "export"
)

// Export writes the environment variables for a given target to the provided
// file in the specified format. If file is nil, output goes to stdout.
func (s *Store) Export(target string, format ExportFormat, file *os.File) error {
	envs, err := s.Load(target)
	if err != nil {
		return fmt.Errorf("load target %q: %w", target, err)
	}

	if file == nil {
		file = os.Stdout
	}

	switch format {
	case FormatDotenv:
		return exportDotenv(envs, file)
	case FormatJSON:
		return exportJSON(envs, file)
	case FormatExport:
		return exportShell(envs, file)
	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

func exportDotenv(envs map[string]string, f *os.File) error {
	for k, v := range envs {
		_, err := fmt.Fprintf(f, "%s=%s\n", k, quoteValue(v))
		if err != nil {
			return err
		}
	}
	return nil
}

func exportJSON(envs map[string]string, f *os.File) error {
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(envs)
}

func exportShell(envs map[string]string, f *os.File) error {
	for k, v := range envs {
		_, err := fmt.Fprintf(f, "export %s=%s\n", k, quoteValue(v))
		if err != nil {
			return err
		}
	}
	return nil
}

// quoteValue wraps a value in single quotes if it contains spaces or special
// characters, escaping any existing single quotes.
func quoteValue(v string) string {
	if strings.ContainsAny(v, " \t\n\r$`\\") {
		return "'" + strings.ReplaceAll(v, "'", "'\\'''") + "'"
	}
	return v
}
