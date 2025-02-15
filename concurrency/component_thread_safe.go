package concurrency

import "sync"

type ComponentThreadSafe struct {
	Data  map[string]string
	Mutex sync.Mutex
}

func NewComponentThreadSafe() *ComponentThreadSafe {
	return &ComponentThreadSafe{Data: make(map[string]string)}
}

func (c *ComponentThreadSafe) HasKey(key string, requireLock bool) bool {
	if requireLock {
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
	}
	if _, ok := c.Data[key]; ok {
		return true
	} else {
		return false
	}
}

func (c *ComponentThreadSafe) Add(key, value string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if c.HasKey(key, false) {
		return
	}
	c.Data[key] = value
}

func (c *ComponentThreadSafe) Read(key string) string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.Data[key]
}
