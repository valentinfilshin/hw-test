package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type CacheItem struct {
	Key   Key
	Value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	cacheItem := CacheItem{
		Key:   key,
		Value: value,
	}

	if val, ok := l.items[key]; ok {
		val.Value = cacheItem
		l.queue.MoveToFront(val)
		return true
	}

	l.queue.PushFront(cacheItem)
	l.items[key] = l.queue.Front()

	if l.queue.Len() > l.capacity {
		oldest := l.queue.Back()
		delete(l.items, oldest.Value.(CacheItem).Key)
		l.queue.Remove(l.queue.Back())
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if val, ok := l.items[key]; ok {
		l.queue.MoveToFront(val)
		return val.Value.(CacheItem).Value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
