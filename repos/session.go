package repos

import (
	"net"
	"time"
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
	Delete(id TToken, opts ...SessionOption) error
	Get(id TToken, opts ...SessionOption) (*Session, error)
}

type sessionOptions struct {
	scope  string
	expire *time.Duration
	idle   *time.Duration
}

type SessionOption interface {
	apply(*sessionOptions)
}

type funcSessionOption struct {
	f func(*sessionOptions)
}

func (it *funcSessionOption) apply(o *sessionOptions) {
	it.f(o)
}

func newFuncSessionOption(f func(*sessionOptions)) SessionOption {
	return &funcSessionOption{f: f}
}
func getSessionOptions(opts ...SessionOption) *sessionOptions {
	o := &sessionOptions{
		scope: "$global",
	}
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

func SessionExpire(expire time.Duration) SessionOption {
	return newFuncSessionOption(func(o *sessionOptions) {
		o.expire = &expire
	})
}

func SessionScope(scope string) SessionOption {
	return newFuncSessionOption(func(o *sessionOptions) {
		o.scope = scope
	})
}

func SessionMaxIdle(idle time.Duration) SessionOption {
	return newFuncSessionOption(func(o *sessionOptions) {
		if idle != 0 {
			o.idle = &idle
		}
	})
}
