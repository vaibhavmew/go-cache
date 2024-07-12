package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInsertWithoutExpiry(t *testing.T) {
	t.Parallel()

	cache := NewCache()

	cache.Insert("key", "value", time.Time{})

	value := cache.Get("key")
	assert.Equal(t, "value", value)
	assert.Empty(t, cache.data["key"].Expiry)
}

func TestInsertWithExpiry(t *testing.T) {
	t.Parallel()

	cache := NewCache()

	cache.Insert("key", "value", time.Now())

	value := cache.Get("key")
	assert.Equal(t, "value", value)
	assert.NotEmpty(t, cache.data["key"].Expiry)
}

func TestDeleteKey(t *testing.T) {
	t.Parallel()

	cache := NewCache()

	cache.Insert("key-two", "value-two", time.Now())
	cache.Insert("key-three", "value-three", time.Now())

	assert.Equal(t, 2, len(cache.data))
	cache.Delete("key-two")
	assert.Equal(t, 1, len(cache.data))

}

func TestUpdateKeyWithExpiry(t *testing.T) {
	t.Parallel()

	cache := NewCache()
	cache.Insert("key-four", "value-four", time.Time{})

	cache.Update("key-four", "value-five", time.Now())

	data := cache.data["key-four"]
	assert.Equal(t, "value-five", data.Value)
	assert.NotEmpty(t, data.Expiry)

}

func TestUpdateKeyWithoutExpiry(t *testing.T) {
	t.Parallel()

	cache := NewCache()
	now := time.Now()
	cache.Insert("key-five", "value-four", now)

	//only updating the key not the expiry time
	cache.Update("key-five", "value-five", time.Time{})

	data := cache.data["key-five"]
	assert.Equal(t, "value-five", data.Value)
	assert.Equal(t, now, data.Expiry)

}

func TestUpdateOnlyExpiry(t *testing.T) {
	t.Parallel()

	cache := NewCache()
	now := time.Now()
	cache.Insert("key-five", "value-four", time.Time{})

	//only updating the key not the expiry time
	cache.Update("key-five", "", now)

	data := cache.data["key-five"]
	assert.Equal(t, "value-four", data.Value)
	assert.Equal(t, now, data.Expiry)

}

func TestUpdateKeyDoesNotExist(t *testing.T) {
	t.Parallel()

	cache := NewCache()
	now := time.Now()

	cache.Update("key-five", "value-five", now)

	data := cache.data["key-five"]
	assert.Equal(t, "value-five", data.Value)
	assert.Equal(t, now, data.Expiry)
}

func TestKeyNotFound(t *testing.T) {
	t.Parallel()

	cache := NewCache()

	value := cache.Get("key")
	assert.Equal(t, "", value)
}

func TestMonitorWithExpiry(t *testing.T) {
	t.Parallel()

	cache := NewCache()

	cache.Insert("key", "value", time.Now())
	assert.Equal(t, 1, len(cache.data))

	stopCh := make(chan struct{})

	go cache.Monitor(stopCh)

	time.Sleep(time.Second * 2)

	close(stopCh)

	value := cache.Get("key")
	assert.Equal(t, "", value)
	assert.Equal(t, 0, len(cache.data))
}

func TestMonitorWithoutExpiry(t *testing.T) {
	t.Parallel()

	cache := NewCache()

	cache.Insert("key", "value", time.Time{})
	assert.Equal(t, 1, len(cache.data))

	stopCh := make(chan struct{})

	go cache.Monitor(stopCh)

	time.Sleep(time.Second * 2)

	close(stopCh)

	value := cache.Get("key")
	assert.Equal(t, "value", value)
	assert.Equal(t, 1, len(cache.data))
}

func TestGetKeys(t *testing.T) {
	t.Parallel()

	cache := NewCache()

	cache.Insert("key", "value", time.Now())
	cache.Insert("key2", "value", time.Now())
	cache.Insert("key3", "value", time.Now())
	cache.Insert("key4", "value", time.Now())

	keys := cache.List()
	assert.Equal(t, 4, len(keys))

}
