package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiffIdentical(t *testing.T) {
	s := tempStore(t)
	require.NoError(t, s.Set("prod", "KEY1", "val1"))
	require.NoError(t, s.Set("prod", "KEY2", "val2"))
	require.NoError(t, s.Set("staging", "KEY1", "val1"))
	require.NoError(t, s.Set("staging", "KEY2", "val2"))

	result, err := s.Diff("prod", "staging")
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"KEY1", "KEY2"}, result.Identical)
	assert.Empty(t, result.OnlyInA)
	assert.Empty(t, result.OnlyInB)
	assert.Empty(t, result.Changed)
}

func TestDiffChanged(t *testing.T) {
	s := tempStore(t)
	require.NoError(t, s.Set("prod", "DB_URL", "prod-db"))
	require.NoError(t, s.Set("staging", "DB_URL", "staging-db"))

	result, err := s.Diff("prod", "staging")
	require.NoError(t, err)
	assert.Equal(t, [2]string{"prod-db", "staging-db"}, result.Changed["DB_URL"])
	assert.Empty(t, result.OnlyInA)
	assert.Empty(t, result.OnlyInB)
}

func TestDiffOnlyInA(t *testing.T) {
	s := tempStore(t)
	require.NoError(t, s.Set("prod", "PROD_ONLY", "secret"))
	require.NoError(t, s.Set("prod", "SHARED", "value"))
	require.NoError(t, s.Set("staging", "SHARED", "value"))

	result, err := s.Diff("prod", "staging")
	require.NoError(t, err)
	assert.Equal(t, []string{"PROD_ONLY"}, result.OnlyInA)
	assert.Empty(t, result.OnlyInB)
}

func TestDiffOnlyInB(t *testing.T) {
	s := tempStore(t)
	require.NoError(t, s.Set("prod", "SHARED", "value"))
	require.NoError(t, s.Set("staging", "SHARED", "value"))
	require.NoError(t, s.Set("staging", "STAGING_ONLY", "debug"))

	result, err := s.Diff("prod", "staging")
	require.NoError(t, err)
	assert.Empty(t, result.OnlyInA)
	assert.Equal(t, []string{"STAGING_ONLY"}, result.OnlyInB)
}

func TestDiffEmptyTargets(t *testing.T) {
	s := tempStore(t)

	result, err := s.Diff("prod", "staging")
	require.NoError(t, err)
	assert.Empty(t, result.OnlyInA)
	assert.Empty(t, result.OnlyInB)
	assert.Empty(t, result.Changed)
	assert.Empty(t, result.Identical)
}
