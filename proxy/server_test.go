package proxy

import (
	"context"
	"net"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEchoServer(t *testing.T) {
	s := NewServer("0.0.0.0:4321", NewFuncProxy(Echo))
	s.Run(context.TODO())
}

func TestForwardServer(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:4322")
	assert.Nil(t, err)
	s := NewServer("0.0.0.0:4321", NewFuncProxy(Forward(addr)))
	s.Run(context.TODO())
}

func TestBenchRecv(t *testing.T) {
	s := &Statd{}
	go s.Tick()
	ctx := WithStatd(context.TODO(), s)
	NewServer("0.0.0.0:4321", NewFuncProxy(BenchRecv())).Run(ctx)
}

func BenchmarkSend(b *testing.B) {
	runtime.GOMAXPROCS(2)
	s := &Statd{}
	go s.Tick()
	ctx := WithStatd(context.TODO(), s)
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:4323")
	NewServer("0.0.0.0:4321", NewFuncProxy(BenchSend(addr)), ServerWorker(2)).Run(ctx)
}
