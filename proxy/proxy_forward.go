package proxy

import (
	"net"
)

func Forward(addr *net.UDPAddr) ProxyHandler {
	return func(ctx *ProxyContext, pes *ProxyEndpoints) error {
		for {
			select {
			case m := <-pes.Downstream.Rx():
				m.Addr = addr
				pes.Upstream.Tx() <- m
			case m := <-pes.Upstream.Rx():
				pes.Downstream.Tx() <- m
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
