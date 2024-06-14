package cache

import "time"

type CacheItem struct {
	Value     []byte
	CreatedAt time.Time
}

type Cache struct {
	myCahce map[string]CacheItem
}

func NewCache(interval time.Duration) *Cache {

	c := &Cache{myCahce: make(map[string]CacheItem)}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Get(key string) ([]byte, bool) {
	if value, ok := c.myCahce[key]; ok {
		return value.Value, true
	}
	return []byte{}, false
}

func (c *Cache) Set(key string, value []byte) bool {
	c.myCahce[key] = CacheItem{Value: value, CreatedAt: time.Now()}
	return true
}

func (c *Cache) reapLoop(interval time.Duration) {
	for {
		time.Sleep(interval)
		for key := range c.myCahce {
			if time.Now().After(c.myCahce[key].CreatedAt.Add(interval)) {
				delete(c.myCahce, key)
			}
		}
	}
}
