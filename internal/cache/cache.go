package cache

type Cache struct {
	myCahce map[string][]byte
}

func NewCache() *Cache {
	return &Cache{myCahce: make(map[string][]byte)}
}

func (c *Cache) Get(key string) ([]byte, error) {
	if value, ok := c.myCahce[key]; ok {
		return value, nil
	}
	return []byte{}, nil
}

func (c *Cache) Set(key string, value []byte) error {
	c.myCahce[key] = value
	return nil
}
