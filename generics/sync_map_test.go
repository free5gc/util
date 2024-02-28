package generics_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/free5gc/util/generics"
)

func TestSyncMap(t *testing.T) {
	var m generics.SyncMap[int, string]

	// Load and Store
	m.Store(1, "test")
	v, e := m.Load(1)
	assert.True(t, e)
	assert.Equal(t, "test", v)
	_, e = m.Load(2)
	assert.False(t, e)

	// Delete
	m.Store(2, "test2")
	m.Delete(1)
	_, e = m.Load(1)
	assert.False(t, e)
	m.Delete(3)
	v, e = m.Load(2)
	assert.True(t, e)
	assert.Equal(t, "test2", v)

	// LoadOrStore
	v, e = m.LoadOrStore(1, "test")
	assert.False(t, e)
	assert.Equal(t, "test", v)
	v, e = m.Load(1)
	assert.True(t, e)
	assert.Equal(t, "test", v)
	v, e = m.LoadOrStore(2, "testX")
	assert.True(t, e)
	assert.Equal(t, "test2", v)
	v, e = m.Load(2)
	assert.True(t, e)
	assert.Equal(t, "test2", v)

	// LoadAndDelete
	v, e = m.LoadAndDelete(1)
	assert.True(t, e)
	assert.Equal(t, "test", v)
	_, e = m.Load(1)
	assert.False(t, e)
	_, e = m.LoadAndDelete(1)
	assert.False(t, e)

	// Swap
	_, e = m.Swap(1, "test")
	assert.False(t, e)
	v, e = m.Load(1)
	assert.True(t, e)
	assert.Equal(t, "test", v)
	v, e = m.Swap(2, "testX")
	assert.True(t, e)
	assert.Equal(t, "test2", v)
	v, e = m.Load(2)
	assert.True(t, e)
	assert.Equal(t, "testX", v)

	// CompareAndDelete
	m.Delete(1)
	m.Store(2, "test2")
	assert.False(t, m.CompareAndDelete(2, "testX"))
	v, e = m.Load(2)
	assert.True(t, e)
	assert.Equal(t, "test2", v)
	assert.True(t, m.CompareAndDelete(2, "test2"))
	_, e = m.Load(2)
	assert.False(t, e)

	// CompareAndSwap
	m.Store(2, "test2")
	assert.False(t, m.CompareAndSwap(1, "test", "testX"))
	_, e = m.Load(1)
	assert.False(t, e)
	assert.False(t, m.CompareAndSwap(2, "testX", "testY"))
	v, e = m.Load(2)
	assert.True(t, e)
	assert.Equal(t, "test2", v)
	assert.True(t, m.CompareAndSwap(2, "test2", "testY"))
	v, e = m.Load(2)
	assert.True(t, e)
	assert.Equal(t, "testY", v)

	m.Range(func(key int, value string) bool {
		switch key {
		case 2:
			assert.Equal(t, "testY", value)
		default:
			assert.Fail(t, "Invalid key %d", key)
		}
		return true
	})
}
