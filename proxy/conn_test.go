package proxy

import (
	"context"
	"testing"
	"net"

	"github.com/libp2p/go-reuseport"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/ipv4"
)

func TestConn(t *testing.T) {
	plisten, err := reuseport.ListenPacket("udp", "0.0.0.0:0")
	assert.Nil(t, err)
	pconn := ipv4.NewPacketConn(plisten)
	conn := NewConn(pconn)
	ctx := context.TODO()
	conn.Run(ctx)
}

func BenchmarkConnWrite(b *testing.B) {
	l, err := reuseport.ListenPacket("udp4", "0.0.0.0:0")
	if err != nil {
		panic(err)
	}
	pconn := ipv4.NewPacketConn(l)
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:4323")
	conn := NewConn(pconn, ConnWriteQueueLength(1024), ConnBatchSize(32))
	go conn.Run(context.TODO())
	umsg := &UdpMsg{Buffer: []byte{0}, Addr: addr}
	for {
		conn.Tx() <- umsg
	}
}
