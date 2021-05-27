package proxy

import (
	"net"
)

func BenchRecv() ProxyHandler {
	return func(ctx *ProxyContext, pes *ProxyEndpoints) error {
		s := GetStatd(ctx.Context)
		for {
			select {
			case msg := <-pes.Downstream.Rx():
				if s != nil {
					s.DRx.Incr(len(msg.Buffer))
				}
				msg.Drop(ctx.BufferPool)
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func BenchSend(target *net.UDPAddr) ProxyHandler {
	return func(ctx *ProxyContext, pes *ProxyEndpoints) error {
		s := GetStatd(ctx.Context)
		for {
			buffer := ctx.BufferPool.Get().([]byte)
			msg := &Msg{Buffer: buffer, Addr: target}
			select {
			case pes.Upstream.Tx() <- msg:
				if s != nil {
					s.UTx.Incr(1)
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
