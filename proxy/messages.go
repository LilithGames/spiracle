package proxy

import (
	"net"

	"golang.org/x/net/ipv4"
)

type Msg = UdpMsg

type Msgs = []ipv4.Message

type UdpMsg struct {
	Buffer []byte
	Addr   net.Addr
}

func (it *UdpMsg) Drop(pool *Pool) {
	pool.Put(it.Buffer[:cap(it.Buffer)])
	it.Buffer = nil
	it.Addr = nil
}
