package roomproxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/LilithGames/spiracle/config"
	"github.com/LilithGames/spiracle/protocol"
	"github.com/LilithGames/spiracle/protocol/heartbeat"
	"github.com/LilithGames/spiracle/protocol/kcp"
	"github.com/LilithGames/spiracle/protocol/multiplex"
	"github.com/LilithGames/spiracle/proxy"
	"github.com/LilithGames/spiracle/repos"
	"github.com/buraksezer/olric"
)

type RoomProxy struct {
	ctx       context.Context
	name      string
	conf      *config.Config
	session   repos.SessionRepo
	router    repos.RouterRepo
	kcp       protocol.Parser
	multiplex protocol.Parser
	heartbeat protocol.Parser
	expire    time.Duration
}

func NewRoomProxy(ctx context.Context, conf *config.Config, name string, router repos.RouterRepo, db *olric.Olric) (*RoomProxy, error) {
	session, err := repos.NewSessionRepo(fmt.Sprintf("roomproxy.%s", name), db)
	if err != nil {
		return nil, fmt.Errorf("NewRoomProxy pstream NewSessionRepo err: %w", err)
	}
	if err != nil {
		return nil, err
	}
	rp := &RoomProxy{
		ctx:       ctx,
		name:      name,
		conf:      conf,
		session:   session,
		router:    router,
		kcp:       protocol.NewFuncParser(kcp.Parser()),
		multiplex: protocol.NewFuncParser(multiplex.Parser()),
		heartbeat: protocol.NewFuncParser(heartbeat.Parser()),
		expire:    time.Second * 30,
	}
	return rp, nil
}

func (it *RoomProxy) token(buffer []byte) (uint32, error) {
	ch, err := it.multiplex.GetToken(buffer)
	if err != nil {
		return 0, err
	}
	switch ch.(byte) {
	case 0x01:
		token, err := it.kcp.GetToken(buffer[1:])
		if err != nil {
			return 0, err
		}
		return token.(uint32), nil
	case 'x':
		token, err := it.heartbeat.GetToken(buffer[1:])
		if err != nil {
			return 0, err
		}
		return token.(uint32), nil
	default:
		return 0, errors.New("unknown channel")
	}
}

func (it *RoomProxy) droute(m *proxy.UdpMsg) error {
	src := m.Addr
	var dst *net.UDPAddr
	token, err := it.token(m.Buffer)
	if err != nil {
		// log
		return err
	}
	s, err := it.session.Get(token)
	if err != nil {
		// warning if !errors.Is(err, repos.ErrKeyNotFound)
		record, err := it.router.Get(token)
		if err != nil {
			// log
			return err
		}
		dst = record.Addr
		// warning
		it.session.CreateOrUpdate(&repos.Session{Token: token, Src: src, Dst: dst}, repos.Expire(it.expire))
	} else {
		dst = s.Dst
	}
	m.Addr = dst
	return nil
}

func (it *RoomProxy) uroute(m *proxy.UdpMsg) error {
	token, err := it.token(m.Buffer)
	if err != nil {
		return err
	}
	s, err := it.session.Get(token)
	if err != nil {
		return err
	}
	m.Addr = s.Src
	return nil
}

func (it *RoomProxy) Run(ctx *proxy.ProxyContext, pes *proxy.ProxyEndpoints) error {
	for {
		select {
		case m := <-pes.Downstream.Rx():
			err := it.droute(m)
			if err != nil {
				m.Drop(ctx.BufferPool)
				continue
			}
			pes.Upstream.Tx() <- m
		case m := <-pes.Upstream.Rx():
			err := it.uroute(m)
			if err != nil {
				m.Drop(ctx.BufferPool)
				continue
			}
			pes.Downstream.Tx() <- m
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
