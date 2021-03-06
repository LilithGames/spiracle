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
	urx    Traffic
	utx    Traffic
	drx    Traffic
	dtx    Traffic
	urxch  [256]Ch
	utxch  [256]Ch
	drxch  [256]Ch
	dtxch  [256]Ch
	pool   Poolx
	udrop  Traffic
	ddrop  Traffic
}

func (it *Statd) URx() *Traffic {
	if it == nil {
		return nil
	}
	return &it.urx
}

func (it *Statd) UTx() *Traffic {
	if it == nil {
		return nil
	}
	return &it.utx
}

func (it *Statd) DRx() *Traffic {
	if it == nil {
		return nil
	}
	return &it.drx
}

func (it *Statd) DTx() *Traffic {
	if it == nil {
		return nil
	}
	return &it.dtx
}

func (it *Statd) URxch(i int) *Ch {
	if it == nil {
		return nil
	}
	return &it.urxch[i]
}
func (it *Statd) UTxch(i int) *Ch {
	if it == nil {
		return nil
	}
	return &it.utxch[i]
}
func (it *Statd) DRxch(i int) *Ch {
	if it == nil {
		return nil
	}
	return &it.drxch[i]
}
func (it *Statd) DTxch(i int) *Ch {
	if it == nil {
		return nil
	}
	return &it.dtxch[i]
}
func (it *Statd) UDrop() *Traffic {
	if it == nil {
		return nil
	}
	return &it.udrop
}
func (it *Statd) DDrop() *Traffic {
	if it == nil {
		return nil
	}
	return &it.ddrop
}

func (it *Statd) Pool() *Poolx {
	if it == nil {
		return nil
	}
	return &it.pool
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
	traffic := fmt.Sprintf("drx: %s dtx: %s urx: %s utx: %s", it.DRx().String(), it.DTx().String(), it.URx().String(), it.UTx().String())
	lines = append(lines, traffic)
	tx := fmt.Sprintf("dtxch %s utxch %s", it.DTxch(0).String(), it.UTxch(0).String())
	lines = append(lines, tx)
	for i := 0; i < it.Worker; i++ {
		rx := fmt.Sprintf("drxch%d %s urxch%d %s", i, it.DRxch(i).String(), i, it.URxch(i).String())
		lines = append(lines, rx)
	}
	pool := fmt.Sprintf("pool %s", it.Pool().String())
	lines = append(lines, pool)
	drop := fmt.Sprintf("ddrop: %s udrop: %s", it.DDrop().String(), it.UDrop().String())
	lines = append(lines, drop)
	return strings.Join(lines, "\n")
}

type TickHandler func(s *Statd)

func StdoutTickHandler(s *Statd) {
	fmt.Printf("%v\n", s.String())
}

func (it *Statd) Tick(th TickHandler) {
	for range time.Tick(time.Second) {
		if th != nil {
			th(it)
		}
		it.URx().Reset()
		it.UTx().Reset()
		it.DRx().Reset()
		it.DTx().Reset()
		it.UDrop().Reset()
		it.DDrop().Reset()
	}
}

func (it *Traffic) Reset() {
	if it == nil {
		return
	}
	atomic.SwapUint64(&it.pps, uint64(0))
	atomic.SwapUint64(&it.bw, uint64(0))
}

func (it *Traffic) Incr(size int) {
	if it == nil {
		return
	}
	atomic.AddUint64(&it.pps, 1)
	atomic.AddUint64(&it.bw, uint64(size))
}

func (it *Traffic) Add(n int, size int) {
	if it == nil {
		return
	}
	atomic.AddUint64(&it.pps, uint64(n))
	atomic.AddUint64(&it.bw, uint64(size))
}

func (it *Traffic) String() string {
	pps := atomic.LoadUint64(&it.pps)
	bw := atomic.LoadUint64(&it.bw)
	return fmt.Sprintf("%v(%v)", numberToUnit(pps), numberToUnit(bw))
}

func (it *Ch) Set(size int) {
	if it == nil {
		return
	}
	atomic.StoreUint64(&it.size, uint64(size))
}

func (it *Ch) String() string {
	size := atomic.LoadUint64(&it.size)
	return fmt.Sprintf("%v", numberToUnit(size))
}

func (it *Poolx) Get() {
	if it == nil {
		return
	}
	atomic.AddInt64(&it.size, 1)
}

func (it *Poolx) Put() {
	if it == nil {
		return
	}
	atomic.AddInt64(&it.size, -1)
}

func (it *Poolx) String() string {
	size := atomic.LoadInt64(&it.size)
	return fmt.Sprintf("%d", size)
}
