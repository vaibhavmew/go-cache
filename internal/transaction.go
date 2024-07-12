package cache

import "sync"

type Transaction struct {
	mu             sync.Mutex      //also lock the cache mutex
	data           map[string]Data //used for storing temporary data
	isLockAcquired bool            //refers to the cache lock
}

func (c *Cache) Start() (*Transaction, error) {
	var (
		transaction Transaction
		err         error
	)

	c.mu.Lock()             //lock the cache mutex
	c.transaction.mu.Lock() //transaction block starts
	c.transaction.isLockAcquired = true
	//empty the data object

	return &transaction, err
}

// read, write, update, delete

func (c *Cache) Commit() error {
	var err error

	//inserting data into the data field
	for k, v := range c.transaction.data {
		c.data[k] = v
	}

	c.transaction.mu.Unlock() //transaction block complete
	c.mu.Unlock()
	c.transaction.isLockAcquired = false
	//empty the data object

	return err
}

func (c *Cache) Abort() error {

	c.transaction.mu.Unlock() //transaction block complete
	c.mu.Unlock()
	c.transaction.isLockAcquired = false
	//empty the data object

	return nil
}
