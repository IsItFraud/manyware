package main

import (
	"errors"

	"github.com/isitfraud/manyware"
)

func main() {
	// The executor runs last, after any Middleware.
	executor := func(m *numMapCtx) error {
		n := m.Get("key1")
		if n != 1 {
			return errors.New("key1 should have been 1")
		}

		n = m.Get("key2")
		if n != 2 {
			return errors.New("key2 should have been 2")
		}

		return nil
	}

	mw1 := func(next manyware.Executor[*numMapCtx]) manyware.Executor[*numMapCtx] {
		return func(m *numMapCtx) error {
			m.Set("key1", 1)
			return next(m)
		}
	}

	mw2 := func(next manyware.Executor[*numMapCtx]) manyware.Executor[*numMapCtx] {
		return func(m *numMapCtx) error {
			m.Set("key2", 2)
			return next(m)
		}
	}

	e := manyware.Prepare(executor, mw1, mw2)
	m := numMapCtx{
		m: make(map[string]int),
	}

	err := e(&m)
	if err != nil {
		panic(err)
	}
}
