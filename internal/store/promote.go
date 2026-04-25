package store

import "fmt"

// PromoteResult holds the result of a promotion operation for a single key.
type PromoteResult struct {
	Key       string
	Overwrite bool
}

// Promote copies all environment variables from the source target to the
// destination target, optionally overwriting existing keys. It returns a
// summary of which keys were written and whether each was an overwrite.
func Promote(s *Store, src, dst string, overwrite bool) ([]PromoteResult, error) {
	if src == dst {
		return nil, fmt.Errorf("source and destination targets must differ: %q", src)
	}

	srcVars, err := s.Load(src)
	if err != nil {
		return nil, fmt.Errorf("load source %q: %w", src, err)
	}
	if len(srcVars) == 0 {
		return nil, fmt.Errorf("source target %q is empty", src)
	}

	dstVars, err := s.Load(dst)
	if err != nil {
		return nil, fmt.Errorf("load destination %q: %w", dst, err)
	}

	var results []PromoteResult
	for k, v := range srcVars {
		_, exists := dstVars[k]
		if exists && !overwrite {
			continue
		}
		dstVars[k] = v
		results = append(results, PromoteResult{Key: k, Overwrite: exists})
	}

	if err := s.Save(dst, dstVars); err != nil {
		return nil, fmt.Errorf("save destination %q: %w", dst, err)
	}

	return results, nil
}
