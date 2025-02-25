package service

import "sync"

type Cache[T any] interface {
	Set(value T)
	Get() T
	Clear()
}

type cache struct {
	value []string

	lock sync.RWMutex
}

func NewImageCache() Cache[[]string] {
	return &cache{}
}

func (c *cache) Set(value []string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value = value
}

func (c *cache) Get() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	buf := make([]string, len(c.value))
	copy(buf, c.value)

	return buf
}

func (c *cache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value = make([]string, 0)
}
