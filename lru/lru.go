package lru

import "container/list"

type Value interface {
	Len() int
}

type Cache struct {
	maxBytes    int64
	nowBytes    int64
	dLinkedList *list.List
	cache       map[string]*list.Element
	OnEvicted   func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:    maxBytes,
		dLinkedList: list.New(),
		cache:       make(map[string]*list.Element),
		OnEvicted:   onEvicted,
	}
}

// Get find value by key
// this element move to front
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.dLinkedList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest remove last element
// call OnEvicted
func (c *Cache) RemoveOldest() {
	ele := c.dLinkedList.Back()
	if ele != nil {
		c.dLinkedList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.dLinkedList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nowBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.dLinkedList.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nowBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nowBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.dLinkedList.Len()
}
