package proxy

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/net/ipv4"
)

type Conn interface {
	Run(ctx context.Context) error
	Endpoint
}

type conn struct {
	*connOptions
	pconn *ipv4.PacketConn
}

func NewConn(pconn *ipv4.PacketConn, opts ...ConnOption) Conn {
	o := getConnOptions(opts...)
	if o.rch == nil {
		o.rch = make(chan *Msg, o.rlen)
	}
	if o.wch == nil {
		o.wch = make(chan *Msg, o.wlen)
	}
	if o.bufferpool == nil {
		o.bufferpool = NewBufferPool(context.TODO(), o.mtu)
	}

	return &conn{
		connOptions: o,
		pconn:       pconn,
	}
}

func (it *conn) Rx() <-chan *Msg {
	return it.rch
}

func (it *conn) Tx() chan<- *Msg {
	return it.wch
}

func (it *conn) initMsgs(msgs Msgs) {
	for i := range msgs {
		msgs[i].Buffers[0] = it.bufferpool.Get().([]byte)
		msgs[i].Addr = nil
	}
}

func (it *conn) putMsgs(msgs Msgs) {
	for i := range msgs {
		buffer := msgs[i].Buffers[0]
		it.bufferpool.Put(buffer[:cap(buffer)])
		msgs[i].Buffers[0] = nil
	}
}

func (it *conn) newMsgs() Msgs {
	msgs := make(Msgs, it.size, it.size)
	for i := range msgs {
		msgs[i].Buffers = make([][]byte, 1)
	}
	return msgs
}

func (it *conn) Run(ctx context.Context) error {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wg.Add(2)
	go func() {
		defer wg.Done()
		it.RunWriteLoop(ctx)
	}()
	go func() {
		defer wg.Done()
		it.RunReadLoop(ctx)
	}()
	<-ctx.Done()
	return ctx.Err()
}

func (it *conn) RunReadLoop(ctx context.Context) error {
	msgs := it.newMsgs()
	it.initMsgs(msgs)
	for {
		it.pconn.SetReadDeadline(time.Now().Add(time.Second))
		n, err := it.pconn.ReadBatch(msgs, 0)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				continue
			}
			log.Fatalf("conn read loop err: %v\n", err)
		}
		for i := 0; i < n; i++ {
			msg := Msg{Buffer: msgs[i].Buffers[0][:msgs[i].N], Addr: msgs[i].Addr.(*net.UDPAddr)}
			select {
			case it.rch <- &msg:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		it.initMsgs(msgs[:n])
	}
}

func (it *conn) batchwch(ctx context.Context, msgs Msgs) (int, error) {
	select {
	case msg := <-it.wch:
		msgs[0].Buffers[0] = msg.Buffer
		msgs[0].Addr = msg.Addr
	case <-ctx.Done():
		return 0, ctx.Err()
	}
	for i := 1; i < len(msgs); i++ {
		select {
		case msg := <-it.wch:
			msgs[i].Buffers[0] = msg.Buffer
			msgs[i].Addr = msg.Addr
		case <-ctx.Done():
			return i, ctx.Err()
		// case <-time.After(time.Second):
			// return i, nil
		default:
			return i, nil
		}
	}
	return len(msgs), nil
}

func (it *conn) RunWriteLoop(ctx context.Context) error {
	msgs := it.newMsgs()
	for {
		n, err := it.batchwch(ctx, msgs)
		if err != nil {
			return err
		}
		it.pconn.SetWriteDeadline(time.Now().Add(time.Second))
		_, err = it.pconn.WriteBatch(msgs[:n], 0)
		// TODO: handle return len here https://man7.org/linux/man-pages/man2/sendmmsg.2.html
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				continue
			}
			log.Fatalf("conn write loop err: %v\n", err)
		}
		it.putMsgs(msgs[:n])
	}
}

type ConnOption interface {
	apply(*connOptions)
}

type funcConnOption struct {
	f func(*connOptions)
}

func (it *funcConnOption) apply(o *connOptions) {
	it.f(o)
}

func newFuncConnOption(f func(*connOptions)) ConnOption {
	return &funcConnOption{f: f}
}
func getConnOptions(opts ...ConnOption) *connOptions {
	o := &connOptions{
		mtu:  2048,
		size: 32,
		rlen: 1024,
		wlen: 1024,
	}
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

type connOptions struct {
	mtu        int
	size       int
	rlen       int
	wlen       int
	rch        chan *Msg
	wch        chan *Msg
	bufferpool *Pool
}

func ConnMtu(mtu int) ConnOption {
	return newFuncConnOption(func(o *connOptions) {
		o.mtu = mtu
	})
}

func ConnBatchSize(size int) ConnOption {
	return newFuncConnOption(func(o *connOptions) {
		o.size = size
	})
}

func ConnReadQueueLength(rlen int) ConnOption {
	return newFuncConnOption(func(o *connOptions) {
		o.rlen = rlen
	})
}

func ConnWriteQueueLength(wlen int) ConnOption {
	return newFuncConnOption(func(o *connOptions) {
		o.wlen = wlen
	})
}

func ConnReadQueue(rch chan *Msg) ConnOption {
	return newFuncConnOption(func(o *connOptions) {
		o.rch = rch
	})
}

func ConnWriteQueue(wch chan *Msg) ConnOption {
	return newFuncConnOption(func(o *connOptions) {
		o.wch = wch
	})
}

func ConnBufferPool(pool *Pool) ConnOption {
	return newFuncConnOption(func(o *connOptions) {
		o.bufferpool = pool
	})
}
