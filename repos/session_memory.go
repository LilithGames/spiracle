package repos

type memorySessionRepo struct {
	// dict *sync.Map
	dict *ConcurrentMap[*Session]
}

func NewMemorySessionRepo() SessionRepo {
	return &memorySessionRepo{
		dict: NewConcurrentMap[*Session](),
	}
}
func (it *memorySessionRepo) key(scope string, token string) string {
	return scope + ":" + token
}

func (it *memorySessionRepo) Create(session *Session, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(session.Token))
	if o.expire == nil {
		it.dict.Set(key, session, 0)
	} else {
		it.dict.Set(key, session, *o.expire)
	}
	return nil
}
func (it *memorySessionRepo) Update(session *Session, opts ...SessionOption) error {
	panic("not implemented")
}
func (it *memorySessionRepo) CreateOrUpdate(session *Session, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(session.Token))
	if o.expire == nil {
		it.dict.Set(key, session, 0)
	} else {
		it.dict.Set(key, session, *o.expire)
	}
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
	session, ok := it.dict.Get(key)
	if !ok {
		return nil, ErrNotExists
	}
	return session, nil
}
