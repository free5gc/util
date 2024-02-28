package generics

import "sync"

type SyncMap[K comparable, V any] struct {
	// type conversion protection, see atomic.Pointer
	_ [0]*K
	_ [0]*V

	m sync.Map
}

func (m *SyncMap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.m.CompareAndDelete(key, old)
}

func (m *SyncMap[K, V]) CompareAndSwap(key K, old V, newValue V) bool {
	return m.m.CompareAndSwap(key, old, newValue)
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.m.Delete(key)
}

func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	if v, o := m.m.Load(key); o {
		return v.(V), true
	} else {
		return
	}
}

func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	if v, l := m.m.LoadAndDelete(key); l {
		return v.(V), true
	} else {
		return
	}
}

func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	if a, l := m.m.LoadOrStore(key, value); l {
		return a.(V), true
	} else {
		return value, false
	}
}

func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

func (m *SyncMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	if p, l := m.m.Swap(key, value); l {
		return p.(V), true
	} else {
		return
	}
}
