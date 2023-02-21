package roomproxy

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/LilithGames/spiracle/proxy"
	"github.com/LilithGames/spiracle/repos"
	"github.com/stretchr/testify/assert"
)

func echo(t *testing.T, ctx context.Context, hostport string) {
	addr, err := net.ResolveUDPAddr("udp4", hostport)
	assert.Nil(t, err)
	conn, err := net.ListenUDP("udp", addr)
	assert.Nil(t, err)
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			conn.WriteToUDP(buf[:n], addr)
		}
	}
}

func client(t *testing.T) *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:4321")
	assert.Nil(t, err)
	conn, err := net.DialUDP("udp4", nil, addr)
	assert.Nil(t, err)
	return conn
}

func TestRoomProxyReal(t *testing.T) {
	s := &proxy.Statd{}
	go s.Tick(func(s *proxy.Statd){})
	ctx := proxy.WithStatd(context.TODO(), s)
	name := "server1"
	roomproxy, err := NewRoomProxy(ctx, name)
	assert.Nil(t, err)
	target, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:10086")
	for i := uint32(1); i < math.MaxInt16; i++ {
		roomproxy.Routers().Create(&repos.RouterRecord{Token: i, Addr: target}, repos.RouterScope(name))
	}
	proxy.NewServer("0.0.0.0:4321", roomproxy).Run(ctx)
}

func TestRoomProxy(t *testing.T) {
	s := &proxy.Statd{}
	// go s.Tick()
	ctx := proxy.WithStatd(context.TODO(), s)
	go echo(t, ctx, "127.0.0.1:10086")
	name := "server1"
	roomproxy, err := NewRoomProxy(ctx, name)
	assert.Nil(t, err)
	token := uint32(4)
	target, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:10086")
	roomproxy.Routers().Create(&repos.RouterRecord{Token: token, Addr: target}, repos.RouterScope(name))
	go proxy.NewServer("0.0.0.0:4321", roomproxy).Run(ctx)
	time.Sleep(time.Second)

	c := client(t)
	c.Write(append([]byte{'e'}, []byte("hello")...))
	buf := make([]byte, 2048)
	n, err := c.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, "ehello", string(buf[:n]))

	data1 := append([]byte{0x01, 0x04, 0x00, 0x00, 0x00}, []byte("12345678901234567890hello1")...)
	c.Write(data1)
	n, err = c.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, string(data1), string(buf[:n]))

	for i := 0; i < 100; i++ {
		c.Write(data1)
		_, err = c.Read(buf)
		if err != nil {
			panic(err)
		}
	}
}

func TestRoomProxyConcurrency(t *testing.T) {
	scope := "TestRoomProxyConcurrency"
	serverSize := 10
	clientPerServerSize := 10
	packetPerClientSize := 1000

	s := &proxy.Statd{}
	ready := make(chan struct{})
	wg := sync.WaitGroup{}
	go s.Tick(proxy.StdoutTickHandler)
	ctx := proxy.WithStatd(context.TODO(), s)
	roomproxy, err := NewRoomProxy(ctx, scope) //, RoomProxyDebug(true)
	assert.Nil(t, err)
	worker := func(ctx context.Context, token uint32) {
		defer wg.Done()
		<-ready
		c := client(t)
		c.Write(append([]byte{'e'}, []byte("hello")...))
		buf := make([]byte, 2048)
		n, err := c.Read(buf)
		assert.Nil(t, err)
		assert.Equal(t, "ehello", string(buf[:n]))
		println("connectivity check ok.")

		data1 := new(bytes.Buffer)
		binary.Write(data1, binary.LittleEndian, []byte{0x01})
		binary.Write(data1, binary.LittleEndian, token)
		binary.Write(data1, binary.LittleEndian, []byte("12345678901234567890hello1"))
		for i := 0; i < packetPerClientSize; i++ {
			c.Write(data1.Bytes())
			n, err = c.Read(buf)
			assert.Nil(t, err)
			if i == 0 {
				assert.Equal(t, data1.String(), string(buf[:n]))
			}
		}
	}

	for i := 0; i < serverSize; i++ {
		addr := fmt.Sprintf("127.0.0.1:%d", 20100+i)
		target, _ := net.ResolveUDPAddr("udp4", addr)
		go echo(t, ctx, addr)
		for j := 0; j < clientPerServerSize; j++ {
			token := uint32(i*clientPerServerSize + j)
			roomproxy.Routers().Create(&repos.RouterRecord{Token: token, Addr: target}, repos.RouterScope(scope))
			wg.Add(1)
			go worker(ctx, token)
		}
	}

	go proxy.NewServer("0.0.0.0:4321", roomproxy, proxy.ServerWorker(4)).Run(ctx)
	time.Sleep(time.Second)
	close(ready)
	wg.Wait()
}
