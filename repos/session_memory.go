package repos

import (
	"sync"
)

type memorySessionRepo struct {
	dict *sync.Map
}

func NewMemorySessionRepo() SessionRepo {
	return &memorySessionRepo{
		dict: &sync.Map{},
	}
}
func (it *memorySessionRepo) key(scope string, token string) string {
	return scope + ":" + token
}

func (it *memorySessionRepo) Create(session *Session, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(session.Token))
	_, loaded := it.dict.LoadOrStore(key, session)
	if loaded {
		return ErrAlreadyExists
	}
	return nil
}
func (it *memorySessionRepo) Update(session *Session, opts ...SessionOption) error {
	panic("not implemented")
}
func (it *memorySessionRepo) CreateOrUpdate(session *Session, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(session.Token))
	it.dict.Store(key, session)
	return nil
}
func (it *memorySessionRepo) Delete(id TToken, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(id))
	it.dict.Delete(key)
	return nil
}
func (it *memorySessionRepo) Get(id TToken, opts ...SessionOption) (*Session, error) {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(id))
	value, ok := it.dict.Load(key)
	if !ok {
		return nil, ErrNotExists
	}
	return value.(*Session), nil
}

