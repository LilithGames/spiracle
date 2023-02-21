package roomproxy

import (
	"context"
	"log"
	"time"

	"github.com/LilithGames/spiracle/protocol"
	"github.com/LilithGames/spiracle/protocol/heartbeat"
	"github.com/LilithGames/spiracle/protocol/kcp"
	"github.com/LilithGames/spiracle/protocol/multiplex"
	"github.com/LilithGames/spiracle/protocol/skcp"
	"github.com/LilithGames/spiracle/proxy"
	"github.com/LilithGames/spiracle/repos"
)

type RoomProxy struct {
	*roomProxyOptions
	ctx       context.Context
	name      string
	multiplex protocol.MultiplexParser
	kcp       protocol.TokenParser
	skcp      protocol.TokenParser
	heartbeat protocol.TokenParser
}

type roomProxyOptions struct {
	expire  time.Duration
	session repos.SessionRepo
	router  repos.RouterRepo
	debug   bool
}

func (it *RoomProxy) Routers() repos.RouterRepo {
	return it.router
}

func NewRoomProxy(ctx context.Context, name string, opts ...RoomProxyOption) (*RoomProxy, error) {
	o := getRoomProxyOptions(opts...)
	if o.session == nil {
		o.session = repos.NewMemorySessionRepo()
	}
	if o.router == nil {
		o.router = repos.NewMemoryRouterRepo()
	}
	rp := &RoomProxy{
		roomProxyOptions: o,
		ctx:              ctx,
		name:             name,
		multiplex:        protocol.NewFuncMultiplexParser(multiplex.Parser()),
		kcp:              protocol.NewFuncTokenParser(kcp.Parser()),
		skcp:             protocol.NewFuncTokenParser(skcp.Parser()),
		heartbeat:        protocol.NewFuncTokenParser(heartbeat.Parser()),
	}
	return rp, nil
}

func (it *RoomProxy) Run(ctx *proxy.ProxyContext, pes *proxy.ProxyEndpoints) error {
	s := proxy.GetStatd(ctx.Context)
	for {
		select {
		case m := <-pes.Downstream.Rx():
			s.DRx().Incr(len(m.Buffer))
			src := m.Addr
			ph := it.multiplexing(m.Buffer)
			err := ph.d(ph.ch, m)
			if err != nil {
				s.DDrop().Incr(len(m.Buffer))
				m.Drop(ctx.BufferPool)
				if it.debug {
					log.Println("[ERROR] RoomProxy drop packet: ", err)
				}
				continue
			}
			if src == m.Addr {
				pes.Downstream.Tx() <- m
			} else {
				pes.Upstream.Tx() <- m
			}
			s.DTx().Incr(len(m.Buffer))
		case m := <-pes.Upstream.Rx():
			s.URx().Incr(len(m.Buffer))
			ph := it.multiplexing(m.Buffer)
			err := ph.u(ph.ch, m)
			if err != nil {
				s.UDrop().Incr(len(m.Buffer))
				m.Drop(ctx.BufferPool)
				continue
				if it.debug {
					log.Println("[ERROR] RoomProxy drop packet: ", err)
				}
			}
			pes.Downstream.Tx() <- m
			s.UTx().Incr(len(m.Buffer))
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

type RoomProxyOption interface {
	apply(*roomProxyOptions)
}

type funcRoomProxyOption struct {
	f func(*roomProxyOptions)
}

func (it *funcRoomProxyOption) apply(o *roomProxyOptions) {
	it.f(o)
}

func newFuncRoomProxyOption(f func(*roomProxyOptions)) RoomProxyOption {
	return &funcRoomProxyOption{f: f}
}
func getRoomProxyOptions(opts ...RoomProxyOption) *roomProxyOptions {
	o := &roomProxyOptions{
		expire: time.Minute * 30,
		debug:  false,
	}
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

func RoomProxyExpire(expire time.Duration) RoomProxyOption {
	return newFuncRoomProxyOption(func(o *roomProxyOptions) {
		if expire > 0 {
			o.expire = expire
		}
	})
}

func RoomProxySessionRepo(session repos.SessionRepo) RoomProxyOption {
	return newFuncRoomProxyOption(func(o *roomProxyOptions) {
		o.session = session
	})
}

func RoomProxyRouterRepo(router repos.RouterRepo) RoomProxyOption {
	return newFuncRoomProxyOption(func(o *roomProxyOptions) {
		o.router = router
	})
}

func RoomProxyDebug(debug bool) RoomProxyOption {
	return newFuncRoomProxyOption(func(o *roomProxyOptions) {
		o.debug = debug
	})
}
