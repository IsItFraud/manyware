package main

import "sync"

// numMapCtx is an insanely contrived use case to demonstrate custom types in manyware.
type numMapCtx struct {
	mux sync.Mutex
	m   map[string]int
}

func (n *numMapCtx) Get(key string) int {
	n.mux.Lock()
	defer n.mux.Unlock()
	return n.m[key]
}

func (n *numMapCtx) Set(key string, value int) {
	n.mux.Lock()
	defer n.mux.Unlock()
	n.m[key] = value
}
