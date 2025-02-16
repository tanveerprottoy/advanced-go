package concurrency

import "sync"

type ComponentThreadSafe struct {
	data  map[string]string
	mutex sync.Mutex
}

func NewComponentThreadSafe() *ComponentThreadSafe {
	return &ComponentThreadSafe{data: make(map[string]string)}
}

func (c *ComponentThreadSafe) HasKey(key string, requireLock bool) bool {
	if requireLock {
		c.mutex.Lock()
		defer c.mutex.Unlock()
	}
	if _, ok := c.data[key]; ok {
		return true
	} else {
		return false
	}
}

func (c *ComponentThreadSafe) Add(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.HasKey(key, false) {
		return
	}
	c.data[key] = value
}

func (c *ComponentThreadSafe) Read(key string) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.data[key]
}
