package app

import "sync"

type ComponentThreadSafe struct {
	Data  map[string]string
	Mutex sync.Mutex
}

func NewComponentThreadSafe() *ComponentThreadSafe {
	c := new(ComponentThreadSafe)
	c.Data = make(map[string]string)
	return c
}

func (c *ComponentThreadSafe) HasKey(key string) bool {
	if _, ok := c.Data[key]; ok {
		return true
	} else {
		return false
	}
}

func (c *ComponentThreadSafe) Add(key, value string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if c.HasKey(key) {
		return
	}
	c.Data[key] = value
}