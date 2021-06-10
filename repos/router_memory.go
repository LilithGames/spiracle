package repos

import (
	"sync"
	"fmt"
	"strings"
)

type memoryRouterRepo struct {
	dict *sync.Map
}

func NewMemoryRouterRepo() RouterRepo {
	return &memoryRouterRepo{
		dict: &sync.Map{},
	}
}

func (it *memoryRouterRepo) Create(record *RouterRecord, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	key := fmt.Sprintf("%s:%s", o.scope, str(record.Token))
	_, loaded := it.dict.LoadOrStore(key, record)
	if loaded {
		return ErrAlreadyExists
	}
	return nil
}

func (it *memoryRouterRepo) Update(record *RouterRecord, opts ...RouterOption) error {
	panic("not implemented")
}

func (it *memoryRouterRepo) CreateOrUpdate(record *RouterRecord, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	key := fmt.Sprintf("%s:%s", o.scope, str(record.Token))
	it.dict.Store(key, record)
	return nil
}

func (it *memoryRouterRepo) Delete(token TToken, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	key := fmt.Sprintf("%s:%s", o.scope, str(token))
	it.dict.Delete(key)
	return nil
}

func (it *memoryRouterRepo) Get(token TToken, opts ...RouterOption) (*RouterRecord, error) {
	o := getRouterOptions(opts...)
	key := fmt.Sprintf("%s:%s", o.scope, str(token))
	value, ok := it.dict.Load(key)
	if !ok {
		return nil, ErrNotExists
	}
	return value.(*RouterRecord), nil
}

func (it *memoryRouterRepo) List(f func(*RouterRecord) bool, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	it.dict.Range(func(key interface{}, value interface{}) bool {
		if strings.HasPrefix(key.(string), o.scope) {
			return f(value.(*RouterRecord))
		}
		return true
	})
	return nil
}
