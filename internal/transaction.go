package cache

import (
	"errors"
	"fmt"
	"time"
)

type Transaction struct {
	data           map[string]Data //used for storing temporary data
	isLockAcquired bool            //refers to the cache lock
}

func (c *Cache) Start() (*Transaction, error) {

	transaction := Transaction{
		data:           make(map[string]Data),
		isLockAcquired: true,
	}

	c.mu.Lock() //lock the cache mutex
	c.transaction = transaction
	//empty the data object

	return &transaction, nil
}

func (c *Cache) Commit() error {
	var err error

	//inserting data into the data field
	for k, v := range c.transaction.data {
		fmt.Println(k, v)
		c.data[k] = v
	}

	//clear the transaction metadata. use sync.Pool in the future
	c.transaction.data = make(map[string]Data)
	c.transaction.isLockAcquired = false

	c.mu.Unlock()

	return err
}

func (c *Cache) Abort() error {

	//clear the transaction metadata. use sync.Pool in the future
	c.transaction.data = make(map[string]Data)
	c.transaction.isLockAcquired = false

	c.mu.Unlock()

	return nil
}

func (tx *Transaction) Insert(key string, value string, expiry time.Time) error {
	if !tx.isLockAcquired {
		return errors.New("transaction not started")
	}

	data := Data{
		Value:  value,
		Expiry: expiry,
	}

	tx.data[key] = data

	return nil
}
