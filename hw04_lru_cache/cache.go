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
	if item, ok := l.items[key]; ok {
		item.Value = value
		l.queue.MoveToFront(item)
		return true
	}

	l.queue.PushFront(value)
	l.items[key] = l.queue.Front()

	if l.queue.Len() > l.capacity {
		oldest := l.queue.Back()
		l.queue.Remove(l.queue.Back())

		for k, v := range l.items {
			if v == oldest {
				delete(l.items, k)
				break
			}
		}
	}

	return false
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
