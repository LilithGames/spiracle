package repos

import (
	"errors"

	"github.com/buraksezer/olric"
)

type sessionRepo struct {
	db   *olric.Olric
	dmap *olric.DMap
}

func NewSessionRepo(db *olric.Olric) (SessionRepo, error) {
	dmap, err := db.NewDMap("cache.sessions")
	if err != nil {
		return nil, err
	}
	return &sessionRepo{
		dmap: dmap,
		db:   db,
	}, nil
}

func (it *sessionRepo) key(scope string, token string) string {
	return scope + ":" + token
}

func (it *sessionRepo) Create(session *Session, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(session.Token))
	if o.expire == nil {
		return it.dmap.PutIf(key, session, olric.IfNotFound)
	} else {
		return it.dmap.PutIfEx(key, session, *o.expire, olric.IfNotFound)
	}
}

func (it *sessionRepo) Update(session *Session, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(session.Token))
	if o.expire == nil {
		return it.dmap.PutIf(key, session, olric.IfFound)
	} else {
		return it.dmap.PutIfEx(key, session, *o.expire, olric.IfFound)
	}
}

func (it *sessionRepo) CreateOrUpdate(session *Session, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(session.Token))
	if o.expire == nil {
		return it.dmap.Put(key, session)
	} else {
		return it.dmap.PutEx(key, session, *o.expire)
	}
}

func (it *sessionRepo) Delete(id TToken, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(id))
	return it.dmap.Delete(key)
}

func (it *sessionRepo) Get(id TToken, opts ...SessionOption) (*Session, error) {
	o := getSessionOptions(opts...)
	key := it.key(o.scope, str(id))
	s, err := it.dmap.Get(key)
	if err != nil {
		if errors.Is(err, olric.ErrKeyNotFound) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	if o.expire != nil {
		it.dmap.Expire(key, *o.expire)
	}
	return s.(*Session), nil
}
