package repos

import (
	"strconv"
	"sync"
	"time"
)

func str(token TToken) string {
	return strconv.FormatUint(uint64(token), 10)
}

type Item[T any] struct {
    value      T
    expiration int64
}

type ConcurrentMap[T any] struct {
    items map[string]Item[T]
    mutex sync.RWMutex
}

func NewConcurrentMap[T any]() *ConcurrentMap[T] {
    cm := &ConcurrentMap[T]{
        items: make(map[string]Item[T]),
    }
    go cm.expirationCleanupWorker()
    return cm
}

func (cm *ConcurrentMap[T]) Set(key string, value T, duration time.Duration) {
    var expiration int64
    if duration > 0 {
        expiration = time.Now().Add(duration).UnixNano()
    }
    cm.mutex.Lock()
    defer cm.mutex.Unlock()
    cm.items[key] = Item[T]{
        value:      value,
        expiration: expiration,
    }
}

func (cm *ConcurrentMap[T]) Get(key string) (T, bool) {
    cm.mutex.RLock()
    defer cm.mutex.RUnlock()
    item, ok := cm.items[key]
    if !ok {
        return *new(T), false
    }
    if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
        return *new(T), false
    }
    return item.value, true
}

func (cm *ConcurrentMap[T]) Delete(key string) {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()
    delete(cm.items, key)
}

func (cm *ConcurrentMap[T]) expirationCleanupWorker() {
    for {
        time.Sleep(time.Second)
        cm.mutex.Lock()
        now := time.Now().UnixNano()
        for key, item := range cm.items {
            if item.expiration > 0 && now > item.expiration {
                delete(cm.items, key)
            }
        }
        cm.mutex.Unlock()
    }
}

