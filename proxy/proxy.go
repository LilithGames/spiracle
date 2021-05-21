package proxy

import (
	"context"
)

type ProxyContext struct {
	context.Context
	BufferPool *Pool
}

type Proxy interface {
	Run(ctx *ProxyContext, pes *ProxyEndpoints) error
}

type ProxyHandler func(ctx *ProxyContext, pes *ProxyEndpoints) error

type funcProxy struct {
	f    ProxyHandler
}

func NewFuncProxy(f ProxyHandler) Proxy {
	return &funcProxy{
		f:    f,
	}
}

func (it *funcProxy) Run(ctx *ProxyContext, pes *ProxyEndpoints) error {
	return it.f(ctx, pes)
}
