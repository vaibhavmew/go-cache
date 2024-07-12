package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransactionInsert(t *testing.T) {
	t.Parallel()

	cache := NewCache()
	tx, err := cache.Start()
	assert.Empty(t, err)

	assert.Empty(t, tx.Insert("key-1", "value", time.Now()))
	assert.Empty(t, tx.Insert("key-2", "value", time.Now()))
	assert.Empty(t, tx.Insert("key-3", "value", time.Now()))

	cache.Commit()

	assert.Equal(t, 3, len(cache.data)) //also check if the transaction block has been cleaned
}

func TestTransactionWithoutLock(t *testing.T) {
	t.Parallel()

	tx := Transaction{}

	err := tx.Insert("key-1", "value", time.Now())
	assert.NotNil(t, err)
	assert.Equal(t, "transaction not started", err.Error())
}

func TestAbortTransaction(t *testing.T) {
	t.Parallel()

	cache := NewCache()
	tx, err := cache.Start()
	assert.Empty(t, err)

	assert.Empty(t, tx.Insert("key-1", "value", time.Now()))
	assert.Empty(t, tx.Insert("key-2", "value", time.Now()))
	assert.Empty(t, tx.Insert("key-3", "value", time.Now()))

	cache.Abort()

	assert.Equal(t, 0, len(cache.data)) //also check if the transaction block has been cleaned
}
