package store

import (
	"fmt"
	"regexp"
	"strings"
)

// LintResult holds a single lint warning for a key in a target.
type LintResult struct {
	Target string
	Key    string
	Issue  string
}

var validKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Lint inspects all keys in the given target and returns a list of warnings.
// It checks for:
//   - keys that are not UPPER_SNAKE_CASE
//   - keys with leading/trailing underscores
//   - empty values
//   - duplicate keys (case-insensitive conflicts)
func Lint(s *Store, target string) ([]LintResult, error) {
	envs, err := s.Load(target)
	if err != nil {
		return nil, fmt.Errorf("lint: load target %q: %w", target, err)
	}

	var results []LintResult
	seen := make(map[string]string) // normalised -> original

	for key, val := range envs {
		// Check naming convention
		if !validKeyPattern.MatchString(key) {
			results = append(results, LintResult{
				Target: target,
				Key:    key,
				Issue:  "key should match UPPER_SNAKE_CASE (A-Z, 0-9, underscore, must start with a letter)",
			})
		}

		// Check leading/trailing underscores
		if strings.HasPrefix(key, "_") || strings.HasSuffix(key, "_") {
			results = append(results, LintResult{
				Target: target,
				Key:    key,
				Issue:  "key has leading or trailing underscore",
			})
		}

		// Check empty value
		if strings.TrimSpace(val) == "" {
			results = append(results, LintResult{
				Target: target,
				Key:    key,
				Issue:  "value is empty",
			})
		}

		// Check case-insensitive duplicates
		norm := strings.ToUpper(key)
		if orig, exists := seen[norm]; exists && orig != key {
			results = append(results, LintResult{
				Target: target,
				Key:    key,
				Issue:  fmt.Sprintf("key conflicts with %q (case-insensitive duplicate)", orig),
			})
		} else {
			seen[norm] = key
		}
	}

	return results, nil
}
