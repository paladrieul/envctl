package store

import "fmt"

// CopyResult holds the outcome of a copy operation.
type CopyResult struct {
	Copied    int
	Skipped   int
	Overwrite bool
}

// Copy duplicates all entries from srcTarget to dstTarget within the given Store.
// If overwrite is false, existing keys in dstTarget are preserved.
func Copy(s *Store, srcTarget, dstTarget string, overwrite bool) (CopyResult, error) {
	if srcTarget == dstTarget {
		return CopyResult{}, fmt.Errorf("source and destination targets must differ")
	}

	src, err := s.Load(srcTarget)
	if err != nil {
		return CopyResult{}, fmt.Errorf("load source target %q: %w", srcTarget, err)
	}
	if len(src) == 0 {
		return CopyResult{}, fmt.Errorf("source target %q not found or empty", srcTarget)
	}

	dst, err := s.Load(dstTarget)
	if err != nil {
		return CopyResult{}, fmt.Errorf("load destination target %q: %w", dstTarget, err)
	}

	result := CopyResult{Overwrite: overwrite}
	for k, v := range src {
		if _, exists := dst[k]; exists && !overwrite {
			result.Skipped++
			continue
		}
		dst[k] = v
		result.Copied++
	}

	if err := s.Save(dstTarget, dst); err != nil {
		return CopyResult{}, fmt.Errorf("save destination target %q: %w", dstTarget, err)
	}

	return result, nil
}
