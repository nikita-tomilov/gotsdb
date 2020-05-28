package storage

import (
	"sync"
	"time"
)

type item struct {
	lastAccess int64
}

type TTLSet struct {
	m map[uint32]*item
	l sync.Mutex
}

func NewTTLSet(ln int, maxTTL int) (m *TTLSet) {
	m = &TTLSet{m: make(map[uint32]*item, ln)}
	go func() {
		for now := range time.Tick(time.Second) {
			m.l.Lock()
			for k, v := range m.m {
				if now.Unix() - v.lastAccess > int64(maxTTL) {
					delete(m.m, k)
				}
			}
			m.l.Unlock()
		}
	}()
	return
}

func (m *TTLSet) Len() int {
	return len(m.m)
}

func (m *TTLSet) Put(k uint32) {
	m.l.Lock()
	it, ok := m.m[k]
	if !ok {
		it = &item{}
		m.m[k] = it
	}
	it.lastAccess = time.Now().Unix()
	m.l.Unlock()
}

func (m *TTLSet) Contains(k uint32) bool {
	m.l.Lock()
	_ , ok := m.m[k]
	m.l.Unlock()
	return ok
}