package roomproxy

import (
	"errors"
	"net"

	"github.com/LilithGames/spiracle/proxy"
	"github.com/LilithGames/spiracle/repos"
)

var ErrDrop = errors.New("drop")

type MsgHandler func(ch byte, m *proxy.UdpMsg) error

type proxyHandler struct {
	ch byte
	u MsgHandler
	d MsgHandler
}

func (it *RoomProxy) multiplexing(buffer []byte) proxyHandler {
	ch, err := it.multiplex.GetChannel(buffer)
	if err != nil {
		return proxyHandler{ch: byte(0), u: it.drop, d: it.drop}
	}
	switch ch {
	case 0x01:
		return proxyHandler{ch: ch, u: it.ukcp, d: it.dkcp}
	case 'x':
		return proxyHandler{ch: ch, u: it.ukcp, d: it.dkcp}
	case 'e':
		return proxyHandler{ch: ch, u: it.uecho, d: it.decho}
	default:
		return proxyHandler{ch: ch, u: it.drop, d: it.drop}
	}

}

func (it *RoomProxy) kcptoken(ch byte, buffer []byte) (uint32, error) {
	switch ch {
	case 0x01:
		token, err := it.kcp.GetToken(buffer[1:])
		if err != nil {
			return 0, err
		}
		return token, nil
	case 'x':
		token, err := it.heartbeat.GetToken(buffer[1:])
		if err != nil {
			return 0, err
		}
		return token, nil
	default:
		return 0, errors.New("unknown kcp channel")
	}
}

func (it *RoomProxy) dkcp(ch byte, m *proxy.UdpMsg) error {
	token, err := it.kcptoken(ch, m.Buffer)
	if err != nil {
		// log
		return err
	}

	src := m.Addr
	var dst *net.UDPAddr
	s, err := it.session.Get(token, repos.SessionScope(it.name))
	if err != nil {
		// mecris count here
		if errors.Is(err, repos.ErrNotExists) {
			record, err := it.router.Get(token, repos.RouterScope(it.name))
			if err != nil {
				// log
				return err
			}
			dst = record.Addr
			// warning
			it.session.CreateOrUpdate(&repos.Session{Token: token, Src: src, Dst: dst}, repos.SessionScope(it.name), repos.SessionExpire(it.expire))
		} else {
			// warning
			return err
		}
	} else {
		dst = s.Dst
	}
	m.Addr = dst
	return nil
}

func (it *RoomProxy) ukcp(ch byte, m *proxy.UdpMsg) error {
	token, err := it.kcptoken(ch, m.Buffer)
	if err != nil {
		return err
	}
	s, err := it.session.Get(token, repos.SessionScope(it.name))
	if err != nil {
		return err
	}
	m.Addr = s.Src
	return nil
}

func (it *RoomProxy) decho(ch byte, m *proxy.UdpMsg) error {
	return nil
}

func (it *RoomProxy) uecho(ch byte, m *proxy.UdpMsg) error {
	return ErrDrop
}

func (it *RoomProxy) drop(ch byte, m *proxy.UdpMsg) error {
	return ErrDrop
}
