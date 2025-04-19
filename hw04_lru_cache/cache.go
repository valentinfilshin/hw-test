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

func (l *lruCache) Set(key Key, value interface{}) bool {
	beenInCache := false
	_, ok := l.items[key]

	if ok {
		beenInCache = true
	}
	l.queue.PushFront(value)
	frontItemLink := l.queue.Front()
	l.items[key] = frontItemLink

	if l.queue.Len() > l.capacity {
		//delete(l.items, l.backItems[l.queue.Back()])
		l.queue.Remove(l.queue.Back())
	}

	return beenInCache
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	val, ok := l.items[key]

	if ok {
		l.queue.MoveToFront(val)
		return val.Value, true
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
