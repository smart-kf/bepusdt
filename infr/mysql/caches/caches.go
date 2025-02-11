package caches

import "sync"

type Cache[T any] interface {
	Get(key string) ([]T, bool)
	Set(key string, data []T)
}

type MemoryCache[T any] struct {
	data sync.Map
}

func (m *MemoryCache[T]) Get(key string) ([]T, bool) {
	data, ok := m.data.Load(key)
	if !ok {
		return nil, false
	}
	return data.([]T), true
}

func (m *MemoryCache[T]) Set(key string, data []T) {
	m.data.Store(key, data)
}

func NewCache[T any]() Cache[T] {
	return &MemoryCache[T]{}
}
