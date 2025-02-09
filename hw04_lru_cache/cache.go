package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List // list of keys
	items    map[Key]*ListItem
	mu       sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if ok {
		existingCacheItem := item.Value.(*cacheItem)
		existingCacheItem.value = value
		c.queue.MoveToFront(item)
		return true
	}

	if c.queue.Len() == c.capacity {
		lastCacheItem := c.queue.Back().Value.(*cacheItem)
		delete(c.items, lastCacheItem.key)
		c.queue.Remove(c.queue.Back())
	}

	item = c.queue.PushFront(&cacheItem{
		key:   key,
		value: value,
	})
	c.items[key] = item
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)
	existingCacheItem := item.Value.(*cacheItem)
	return existingCacheItem.value, true
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
