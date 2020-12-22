package lru

import "container/list"

//lru 最近最久未使用的淘汰策略

type Cache struct {
	maxBytes  uint64
	usedBytes uint64
	//双向链表用来存数据
	dl *list.List
	//key 是字符串，值是双向链表中对应的
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

//每个链节点中存储的值类型
type entry struct {
	key   string
	value Value
}

//任何实现了此接口的类型都可以作为值
type Value interface {
	Len() int
}

func New(maxBytes uint64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		dl:        list.New(),
		cache:     make(map[string]*list.Element, 0),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.dl.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.dl.Back()
	if ele != nil {
		c.dl.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.usedBytes -= uint64(len(kv.key)) + uint64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.dl.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.usedBytes = c.usedBytes + uint64(value.Len()) - uint64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.dl.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.usedBytes += uint64(len(key)) + uint64(value.Len())
	}
	for c.maxBytes != 0 && c.usedBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.dl.Len()
}
