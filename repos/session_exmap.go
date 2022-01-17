package repos

import (
	"time"

	"github.com/nursik/go-expire-map"
)

type sessionRepoV2 struct {
	exmap *expiremap.ExpireMap

	Expiration time.Duration
	MaxIdle    *time.Duration
}

func NewSessionRepoV2(opts ...SessionOption) (SessionRepo, error) {
	o := getSessionOptions(opts...)
	expire := time.Second * 30
	if o.expire != nil {
		expire = *o.expire 
	}
	return &sessionRepoV2{
		exmap: expiremap.New(),
		Expiration: expire,
		MaxIdle:    o.idle,
	}, nil
}

func (it *sessionRepoV2) Create(session *Session, opts ...SessionOption) error {
	panic("not implemented")
	return nil
}

func (it *sessionRepoV2) Update(session *Session, opts ...SessionOption) error {
	panic("not implemented")
	return nil
}

func (it *sessionRepoV2) CreateOrUpdate(session *Session, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(session.Token))
	expire := it.Expiration
	if o.expire != nil {
        expire = *o.expire
	}
	it.exmap.Set(key, session, expire)
	return nil
}

func (it *sessionRepoV2) Delete(id TToken, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(id))
	it.exmap.Delete(key)
	return nil
}

func (it *sessionRepoV2) Get(id TToken, opts ...SessionOption) (*Session, error) {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(id))
	v, ok := it.exmap.Get(key)
	if !ok {
		return nil, ErrNotExists
	}
	if it.MaxIdle != nil {
		it.exmap.SetTTL(key, *it.MaxIdle)
	}
	return v.(*Session), nil

}

func (it *sessionRepoV2) key(scope string, token string) string {
	return scope + ":" + token
}

