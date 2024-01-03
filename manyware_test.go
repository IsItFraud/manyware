package manyware_test

import (
	"sync"
	"testing"

	"github.com/isitfraud/manyware"
	"github.com/stretchr/testify/require"
)

func TestPrepare(t *testing.T) {
	// The executor runs last, after any Middleware.
	executor := func(m *sync.Map) error {
		value, ok := m.Load("calls")
		require.True(t, ok)
		require.IsType(t, []string{}, value)
		calls := value.([]string)
		calls = append(calls, "executor")
		m.Store("calls", calls)
		return nil
	}

	mw1 := func(next manyware.Executor[*sync.Map]) manyware.Executor[*sync.Map] {
		return func(m *sync.Map) error {
			// "calls" should not be set yet because this is the first call in the chain.
			_, ok := m.Load("calls")
			require.False(t, ok)

			m.Store("calls", []string{"mw1"})
			return next(m)
		}
	}

	mw2 := func(next manyware.Executor[*sync.Map]) manyware.Executor[*sync.Map] {
		return func(m *sync.Map) error {
			value, ok := m.Load("calls")
			require.True(t, ok)
			require.IsType(t, []string{}, value)
			calls := value.([]string)
			calls = append(calls, "mw2")
			m.Store("calls", calls)
			return next(m)
		}
	}

	s := manyware.Prepare(executor, mw1, mw2)
	m := new(sync.Map)
	err := s(m)
	require.NoError(t, err)

	value, ok := m.Load("calls")
	require.True(t, ok)
	require.IsType(t, []string{}, value)
	calls := value.([]string)
	require.Equal(t, []string{"mw1", "mw2", "executor"}, calls)
}
