package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ImportFormat represents the format of the import source.
type ImportFormat string

const (
	FormatDotenv ImportFormat = "dotenv"
	FormatJSON   ImportFormat = "json"
)

// ImportResult holds the results of an import operation.
type ImportResult struct {
	Imported  int
	Skipped   int
	Overwrite bool
}

// Import reads key-value pairs from r in the given format and stores them
// under target. If overwrite is false, existing keys are skipped.
func Import(s *Store, target string, r io.Reader, format ImportFormat, overwrite bool) (ImportResult, error) {
	var pairs map[string]string
	var err error

	switch format {
	case FormatDotenv:
		pairs, err = parseDotenv(r)
	case FormatJSON:
		pairs, err = parseJSON(r)
	default:
		return ImportResult{}, fmt.Errorf("unsupported import format: %s", format)
	}
	if err != nil {
		return ImportResult{}, fmt.Errorf("parse error: %w", err)
	}

	existing, err := s.Load(target)
	if err != nil {
		return ImportResult{}, err
	}

	result := ImportResult{Overwrite: overwrite}
	for k, v := range pairs {
		if _, ok := existing[k]; ok && !overwrite {
			result.Skipped++
			continue
		}
		existing[k] = v
		result.Imported++
	}

	if err := s.Save(target, existing); err != nil {
		return ImportResult{}, err
	}
	return result, nil
}

func parseDotenv(r io.Reader) (map[string]string, error) {
	pairs := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.Trim(strings.TrimSpace(line[idx+1:]), `"`)
		pairs[key] = val
	}
	return pairs, scanner.Err()
}

func parseJSON(r io.Reader) (map[string]string, error) {
	pairs := make(map[string]string)
	if err := json.NewDecoder(r).Decode(&pairs); err != nil {
		return nil, err
	}
	return pairs, nil
}
