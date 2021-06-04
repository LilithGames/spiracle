package repos

import (
	"errors"
	"net"
)

type TToken = uint32

var ErrAlreadyExists = errors.New("elements already exists")
var ErrNotExists = errors.New("elements not exists")

type RouterRecord struct {
	Token    TToken
	Addr     *net.UDPAddr
	RoomId   string
	PlayerId string
}

type RouterRepo interface {
	Create(record *RouterRecord, opts ...RouterOption) error
	Update(record *RouterRecord, opts ...RouterOption) error
	CreateOrUpdate(record *RouterRecord, opts ...RouterOption) error
	Delete(token TToken, opts ...RouterOption) error
	Get(token TToken, opts ...RouterOption) (*RouterRecord, error)
	List(f func(*RouterRecord) bool, opts ...RouterOption) error
}

type routerOptions struct {
	scope string
}

type RouterOption interface {
	apply(*routerOptions)
}

type funcRouterOption struct {
	f func(*routerOptions)
}

func (it *funcRouterOption) apply(o *routerOptions) {
	it.f(o)
}

func newFuncRouterOption(f func(*routerOptions)) RouterOption {
	return &funcRouterOption{f: f}
}
func getRouterOptions(opts ...RouterOption) *routerOptions {
	o := &routerOptions{
		scope: "$global",
	}
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

func RouterScope(scope string) RouterOption {
	return newFuncRouterOption(func(o *routerOptions) {
		o.scope = scope
	})
}
