package proxy

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

type Traffic struct {
	pps uint64
	bw  uint64
}

type Ch struct {
	size uint64
}

type Poolx struct {
	size int64
}

type Statd struct {
	Name   string
	Worker int
	URx    Traffic
	UTx    Traffic
	DRx    Traffic
	DTx    Traffic
	URxch  [256]Ch
	UTxch  [256]Ch
	DRxch  [256]Ch
	DTxch  [256]Ch
	Pool   Poolx
}

func WithStatd(ctx context.Context, s *Statd) context.Context {
	return context.WithValue(ctx, "spiracle.proxy.statd", s)
}
func GetStatd(ctx context.Context) *Statd {
	s, ok := ctx.Value("spiracle.proxy.statd").(*Statd)
	if !ok {
		return nil
	}
	return s
}

func (it *Statd) String() string {
	lines := []string{}
	traffic := fmt.Sprintf("drx: %s dtx: %s urx: %s utx: %s", it.DRx.String(), it.DTx.String(), it.URx.String(), it.UTx.String())
	lines = append(lines, traffic)
	tx := fmt.Sprintf("dtxch %s utxch %s", (&it.DTxch[0]).String(), (&it.UTxch[0]).String())
	lines = append(lines, tx)
	for i := 0; i < it.Worker; i++ {
		rx := fmt.Sprintf("drxch%d %s urxch%d %s", i, (&it.DRxch[i]).String(), i, (&it.URxch[i]).String())
		lines = append(lines, rx)
	}
	pool := fmt.Sprintf("pool %s", it.Pool.String())
	lines = append(lines, pool)
	return strings.Join(lines, "\n")
}

func (it *Statd) Tick() {
	for range time.Tick(time.Second) {
		fmt.Printf("%v\n", it.String())
		it.URx.Reset()
		it.UTx.Reset()
		it.DRx.Reset()
		it.DTx.Reset()
	}
}

func (it *Traffic) Reset() {
	atomic.SwapUint64(&it.pps, uint64(0))
	atomic.SwapUint64(&it.bw, uint64(0))
}

func (it *Traffic) Incr(size int) {
	atomic.AddUint64(&it.pps, 1)
	atomic.AddUint64(&it.bw, uint64(size))
}

func (it *Traffic) Add(n int, size int) {
	atomic.AddUint64(&it.pps, uint64(n))
	atomic.AddUint64(&it.bw, uint64(size))
}

func (it *Traffic) String() string {
	pps := atomic.LoadUint64(&it.pps)
	bw := atomic.LoadUint64(&it.bw)
	return fmt.Sprintf("%v(%v)", numberToUnit(pps), numberToUnit(bw))
}

func (it *Ch) Set(size int) {
	atomic.StoreUint64(&it.size, uint64(size))
}

func (it *Ch) String() string {
	size := atomic.LoadUint64(&it.size)
	return fmt.Sprintf("%v", numberToUnit(size))
}

func (it *Poolx) Get() {
	atomic.AddInt64(&it.size, 1)
}

func (it *Poolx) Put() {
	atomic.AddInt64(&it.size, -1)
}

func (it *Poolx) String() string {
	size := atomic.LoadInt64(&it.size)
	return fmt.Sprintf("%d", size)
}
