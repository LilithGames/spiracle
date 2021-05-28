package proxy

import (
	"context"
	"sync"
)

type Pool = DefaultPool

type DefaultPool struct {
	sync.Pool
	*Statd
}

func NewBufferPool(ctx context.Context, mtu int) *Pool {
	return &DefaultPool{
		Pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, mtu, mtu)
			},
		},
		Statd: GetStatd(ctx),
	}
}

func (it *DefaultPool) Get() interface{} {
	it.Statd.Pool().Get()
	return it.Pool.Get()
}

func (it *DefaultPool) Put(x interface{}) {
	it.Statd.Pool().Put()
	it.Pool.Put(x)
}
