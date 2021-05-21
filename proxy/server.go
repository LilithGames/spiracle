package proxy

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-reuseport"
	"golang.org/x/net/ipv4"
)

type Server struct {
	*serverOptions
	addr  string
	proxy Proxy
}

type serverOptions struct {
	n      int
	mtu    int
	batch  int
	rxsize int
	txsize int
	upaddr string
	pool   *Pool
}

func NewServer(addr string, proxy Proxy, opts ...ServerOption) *Server {
	o := getServerOptions(opts...)
	return &Server{
		serverOptions: o,
		addr:          addr,
		proxy:         proxy,
	}
}

func (it *Server) Run(ctx context.Context) error {
	pool := it.pool
	if pool == nil {
		pool = NewBufferPool(ctx, it.mtu)
	}
	pctx := &ProxyContext{Context: ctx, BufferPool: pool}
	utx := make(chan *Msg, it.txsize)
	defer close(utx)
	dtx := make(chan *Msg, it.txsize)
	defer close(dtx)
	urx := make([]chan *Msg, it.n)
	for i := 0; i < it.n; i++ {
		urx[i] = make(chan *Msg, it.rxsize)
		defer close(urx[i])
	}
	drx := make([]chan *Msg, it.n)
	for i := 0; i < it.n; i++ {
		drx[i] = make(chan *Msg, it.rxsize)
		defer close(drx[i])
	}
	wg := sync.WaitGroup{}
	defer wg.Wait()

	// start
	for i := 0; i < it.n; i++ {
		upstream := &DuplexEndpoint{Receive: urx[i], Transmit: utx}
		downstream := &DuplexEndpoint{Receive: drx[i], Transmit: dtx}
		ulisten, err := it.newListen(it.upaddr)
		if err != nil {
			log.Fatal("create upstream listen err: ", err)
		}
		uconn := NewConn(ulisten, ConnBufferPool(pool), ConnBatchSize(it.batch), ConnReadQueue(upstream.Receive), ConnWriteQueue(upstream.Transmit))
		wg.Add(1)
		go func() {
			defer wg.Done()
			uconn.Run(ctx)
		}()
		dlisten, err := it.newListen(it.addr)
		if err != nil {
			log.Fatal("create downstream listen err: ", err)
		}
		dconn := NewConn(dlisten, ConnBufferPool(pool), ConnBatchSize(it.batch), ConnReadQueue(downstream.Receive), ConnWriteQueue(downstream.Transmit))
		wg.Add(1)
		go func() {
			defer wg.Done()
			dconn.Run(ctx)
		}()

		//run proxy
		wg.Add(1)
		go func() {
			defer wg.Done()
			it.proxy.Run(pctx, &ProxyEndpoints{Upstream: upstream, Downstream: downstream})
		}()
	}
	report := func(ctx context.Context) {
		s := GetStatd(ctx)
		if s == nil {
			return
		}
		s.Worker = it.n
		for {
			for i := 0; i < it.n; i++ {
				(&s.URxch[i]).Set(len(urx[i]))
				(&s.DRxch[i]).Set(len(drx[i]))
			}
			(&s.UTxch[0]).Set(len(utx))
			(&s.DTxch[0]).Set(len(dtx))
			if ctx.Err() != nil {
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
	go report(ctx)

	<-ctx.Done()
	return ctx.Err()
}

func (it *Server) newListen(addr string) (*ipv4.PacketConn, error) {
	l, err := reuseport.ListenPacket("udp4", addr)
	if err != nil {
		return nil, fmt.Errorf("ListenPacket err: %w", err)
	}
	pconn := ipv4.NewPacketConn(l)
	return pconn, nil
}

type ServerOption interface {
	apply(*serverOptions)
}

type funcServerOption struct {
	f func(*serverOptions)
}

func (it *funcServerOption) apply(o *serverOptions) {
	it.f(o)
}

func newFuncServerOption(f func(*serverOptions)) ServerOption {
	return &funcServerOption{f: f}
}
func getServerOptions(opts ...ServerOption) *serverOptions {
	o := &serverOptions{
		n:      1,
		mtu:    2048,
		batch:  32,
		rxsize: 1024,
		txsize: 1024,
		upaddr: "0.0.0.0:0",
	}
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

func ServerWorker(worker int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.n = worker
	})
}

func ServerMtu(mtu int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.mtu = mtu
	})
}

func ServerBatch(batch int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.batch = batch
	})
}

func ServerUpaddr(upaddr string) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.upaddr = upaddr
	})
}

func ServerRxsize(rxsize int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.rxsize = rxsize
	})
}

func ServerTxsize(txsize int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.txsize = txsize
	})
}

func ServerBufferPool(pool *Pool) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.pool = pool
	})
}
