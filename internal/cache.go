package cache

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"time"
)

type Cache struct {
	mu   sync.Mutex //protects the data map
	data map[string]Data
}

type Data struct {
	Value  string
	Expiry time.Time
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]Data),
	}
}

func (c *Cache) Insert(key string, value string, expiry time.Time) {
	c.mu.Lock()

	data := Data{
		Value:  value,
		Expiry: expiry,
	}

	c.data[key] = data
	c.mu.Unlock()
}

func (c *Cache) Get(key string) string {
	c.mu.Lock()
	val := c.data[key].Value
	c.mu.Unlock()

	return val
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

func (c *Cache) Update(key string, value string, expiry time.Time) {

	c.mu.Lock()

	if expiry.IsZero() {
		expiry = c.data[key].Expiry
	}

	data := Data{
		Value:  value,
		Expiry: expiry,
	}

	c.data[key] = data
	c.mu.Unlock()
}

func (c *Cache) List() []string {
	c.mu.Lock()

	var keys []string

	for key := range c.data {
		keys = append(keys, key)
	}

	c.mu.Unlock()

	return keys
}

func (c *Cache) Monitor(stopCh chan struct{}) {

	for {
		select {
		case <-stopCh:
			fmt.Println("exiting the monitor")
			return
		default:
			time.Sleep(1 * time.Duration(time.Second))

			c.mu.Lock()

			for k, v := range c.data {
				if v.Expiry.Before(time.Now()) {
					delete(c.data, k)
					fmt.Println("deleted the key: ", k)
				}
			}

			c.mu.Unlock()
		}
	}

}

func (c *Cache) Flush() error {
	file, err := os.Create("keys.gob")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)

	err = encoder.Encode(c.data)
	if err != nil {
		return err
	}

	return err
}

func (c *Cache) LoadKeys() error {
	file, err := os.Open("keys.gob")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	var data map[string]Data
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}

	c.data = data
	fmt.Println("loaded ", len(c.data), " keys")
	return err
}
