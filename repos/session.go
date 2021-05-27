package repos

import (
	"net"
	"time"

	"github.com/buraksezer/olric"
)

type Session struct {
	Token TToken
	Src   *net.UDPAddr
	Dst   *net.UDPAddr
}

type SessionRepo interface {
	Create(session *Session, opts ...SessionOption) error
	Update(session *Session, opts ...SessionOption) error
	CreateOrUpdate(session *Session, opts ...SessionOption) error
	Delete(id TToken) error
	Get(id TToken, opts ...SessionOption) (*Session, error)
}

type sessionRepo struct {
	name string
	db   *olric.Olric
	dmap *olric.DMap
}

func NewSessionRepo(scope string, db *olric.Olric) (SessionRepo, error) {
	if scope == "" {
		scope = "$global"
	}
	name := "cache." + scope + ".session"
	dmap, err := db.NewDMap(name)
	if err != nil {
		return nil, err
	}
	return &sessionRepo{
		name: name,
		dmap: dmap,
		db:   db,
	}, nil
}

func (it *sessionRepo) Create(session *Session, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	if o.expire == nil {
		return it.dmap.PutIf(str(session.Token), session, olric.IfNotFound)
	} else {
		return it.dmap.PutIfEx(str(session.Token), session, *o.expire, olric.IfNotFound)
	}
}

func (it *sessionRepo) Update(session *Session, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	if o.expire == nil {
		return it.dmap.PutIf(str(session.Token), session, olric.IfFound)
	} else {
		return it.dmap.PutIfEx(str(session.Token), session, *o.expire, olric.IfFound)
	}
}

func (it *sessionRepo) CreateOrUpdate(session *Session, opts ...SessionOption) error {
	o := getSessionOptions(opts...)
	if o.expire == nil {
		return it.dmap.Put(str(session.Token), session)
	} else {
		return it.dmap.PutEx(str(session.Token), session, *o.expire)
	}
}

func (it *sessionRepo) Delete(id TToken) error {
	return it.dmap.Delete(str(id))
}

func (it *sessionRepo) Get(id TToken, opts ...SessionOption) (*Session, error) {
	token := str(id)
	s, err := it.dmap.Get(token)
	if err != nil {
		return nil, err
	}
	o := getSessionOptions(opts...)
	if o.expire != nil {
		it.dmap.Expire(token, *o.expire)
	}
	return s.(*Session), nil
}

type sessionOptions struct {
	expire *time.Duration
}

type SessionOption interface {
	apply(*sessionOptions)
}

type funcSessionCreateOption struct {
	f func(*sessionOptions)
}

func (it *funcSessionCreateOption) apply(o *sessionOptions) {
	it.f(o)
}

func newFuncSessionCreateOption(f func(*sessionOptions)) SessionOption {
	return &funcSessionCreateOption{f: f}
}
func getSessionOptions(opts ...SessionOption) *sessionOptions {
	o := &sessionOptions{}
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

func Expire(expire time.Duration) SessionOption {
	return newFuncSessionCreateOption(func(o *sessionOptions) {
		o.expire = &expire
	})
}
