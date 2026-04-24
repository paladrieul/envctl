package store

import "sort"

// DiffResult holds the differences between two targets.
type DiffResult struct {
	OnlyInA  []string          // keys present only in target A
	OnlyInB  []string          // keys present only in target B
	Changed  map[string][2]string // keys present in both but with different values
	Identical []string          // keys present in both with the same value
}

// Diff compares the environment variables of two targets and returns the differences.
// Values are compared in plaintext; encrypted values must be decrypted before calling Diff.
func (s *Store) Diff(targetA, targetB string) (*DiffResult, error) {
	envA, err := s.Load(targetA)
	if err != nil {
		return nil, err
	}
	envB, err := s.Load(targetB)
	if err != nil {
		return nil, err
	}

	result := &DiffResult{
		Changed: make(map[string][2]string),
	}

	for k, vA := range envA {
		if vB, ok := envB[k]; ok {
			if vA == vB {
				result.Identical = append(result.Identical, k)
			} else {
				result.Changed[k] = [2]string{vA, vB}
			}
		} else {
			result.OnlyInA = append(result.OnlyInA, k)
		}
	}

	for k := range envB {
		if _, ok := envA[k]; !ok {
			result.OnlyInB = append(result.OnlyInB, k)
		}
	}

	sort.Strings(result.OnlyInA)
	sort.Strings(result.OnlyInB)
	sort.Strings(result.Identical)

	return result, nil
}
