package mmap

import "sync"

type MMap struct {
	m   map[interface{}]interface{}
	mut sync.RWMutex
}

func New() *MMap {
	return &MMap{
		m: make(map[interface{}]interface{}),
	}
}

func (m *MMap) Set(key, value interface{}) {
	m.mut.Lock()
	m.m[key] = value
	m.mut.Unlock()
}

func (m *MMap) Has(key interface{}) bool {
	m.mut.RLock()
	_, has := m.m[key]
	m.mut.RUnlock()
	return has
}

func (m *MMap) Get(key interface{}) interface{} {
	m.mut.RLock()
	value := m.m[key]
	m.mut.RUnlock()
	return value
}

func (m *MMap) Delete(key interface{}) {
	m.mut.Lock()
	delete(m.m, key)
	m.mut.Unlock()
}

func (m *MMap) Keys() interface{} {
	keys := make([]interface{}, len(m.m))
	i := 0
	m.mut.Lock()
	for key := range m.m {
		keys[i] = key
		i++
	}
	m.mut.Unlock()
	return keys
}
