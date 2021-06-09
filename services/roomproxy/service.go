package roomproxy

import (
	"context"
	"fmt"
	"time"

	"github.com/LilithGames/spiracle/infra/db"
	"github.com/LilithGames/spiracle/protocol"
	"github.com/LilithGames/spiracle/protocol/heartbeat"
	"github.com/LilithGames/spiracle/protocol/kcp"
	"github.com/LilithGames/spiracle/protocol/multiplex"
	"github.com/LilithGames/spiracle/proxy"
	"github.com/LilithGames/spiracle/repos"
	"github.com/buraksezer/olric"
)

type RoomProxy struct {
	*roomProxyOptions
	ctx       context.Context
	name      string
	kcp       protocol.Parser
	multiplex protocol.Parser
	heartbeat protocol.Parser
}

type roomProxyOptions struct {
	expire  time.Duration
	session repos.SessionRepo
	router  repos.RouterRepo
	db      *olric.Olric
}

func (it *RoomProxy) Routers() repos.RouterRepo {
	return it.router
}

func NewRoomProxy(ctx context.Context, name string, opts ...RoomProxyOption) (*RoomProxy, error) {
	o := getRoomProxyOptions(opts...)
	var err error
	if o.db == nil {
		o.db, err = db.ProvideServer(ctx, db.ServerLocalConfig())
		if err != nil {
			return nil, fmt.Errorf("NewRoomProxy db.ProvideLocal err: %w", err)
		}
	}
	if o.session == nil {
		o.session, err = repos.NewSessionRepo(fmt.Sprintf("roomproxy.%s", name), o.db)
		if err != nil {
			return nil, fmt.Errorf("NewRoomProxy NewSessionRepo err: %w", err)
		}
	}
	if o.router == nil {
		o.router, err = repos.NewRouterRepo(o.db)
		if err != nil {
			return nil, fmt.Errorf("NewRoomProxy NewRouterRepo err: %w", err)
		}
	}
	rp := &RoomProxy{
		roomProxyOptions: o,
		ctx:              ctx,
		name:             name,
		kcp:              protocol.NewFuncParser(kcp.Parser()),
		multiplex:        protocol.NewFuncParser(multiplex.Parser()),
		heartbeat:        protocol.NewFuncParser(heartbeat.Parser()),
	}
	return rp, nil
}

func (it *RoomProxy) Run(ctx *proxy.ProxyContext, pes *proxy.ProxyEndpoints) error {
	s := proxy.GetStatd(ctx.Context)
	for {
		select {
		case m := <-pes.Downstream.Rx():
			s.DRx().Incr(len(m.Buffer))
			ph := it.multiplexing(m.Buffer)
			err := ph.d(ph.ch, m)
			if err != nil {
				s.DDrop().Incr(len(m.Buffer))
				m.Drop(ctx.BufferPool)
				continue
			}
			pes.Upstream.Tx() <- m
			s.DTx().Incr(len(m.Buffer))
		case m := <-pes.Upstream.Rx():
			s.URx().Incr(len(m.Buffer))
			ph := it.multiplexing(m.Buffer)
			err := ph.u(ph.ch, m)
			if err != nil {
				s.UDrop().Incr(len(m.Buffer))
				m.Drop(ctx.BufferPool)
				continue
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
		expire: time.Second * 30,
	}
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

func RoomProxyExpire(expire time.Duration) RoomProxyOption {
	return newFuncRoomProxyOption(func(o *roomProxyOptions) {
		o.expire = expire
	})
}

func RoomProxyDb(db *olric.Olric) RoomProxyOption {
	return newFuncRoomProxyOption(func(o *roomProxyOptions) {
		o.db = db
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
