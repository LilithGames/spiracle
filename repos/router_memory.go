package repos

import "sync"

type memoryRouterRepo struct {
	dict *sync.Map
}

func NewMemoryRouterRepo() (RouterRepo, error) {
	return &memoryRouterRepo{
		dict: &sync.Map{},
	}, nil
}

func (it *memoryRouterRepo) Create(record *RouterRecord, opts ...RouterOption) error {
	_, loaded := it.dict.LoadOrStore(record.Token, record)
	if loaded {
		return ErrAlreadyExists
	}
	return nil
}

func (it *memoryRouterRepo) Update(record *RouterRecord, opts ...RouterOption) error {
	panic("not implemented")
}

func (it *memoryRouterRepo) CreateOrUpdate(record *RouterRecord, opts ...RouterOption) error {
	it.dict.Store(record.Token, record)
	return nil
}

func (it *memoryRouterRepo) Delete(token TToken, opts ...RouterOption) error {
	it.dict.Delete(token)
	return nil
}

func (it *memoryRouterRepo) Get(token TToken, opts ...RouterOption) (*RouterRecord, error) {
	value, ok := it.dict.Load(token)
	if !ok {
		return nil, ErrNotExists
	}
	return value.(*RouterRecord), nil
}

func (it *memoryRouterRepo) List(f func(*RouterRecord) bool, opts ...RouterOption) error {
	it.dict.Range(func(key interface{}, value interface{}) bool {
		return f(value.(*RouterRecord))
	})
	return nil
}
