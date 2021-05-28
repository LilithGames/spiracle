package roomproxy

import (
	"context"
	"testing"
	"net"

	"github.com/LilithGames/spiracle/proxy"
	"github.com/LilithGames/spiracle/repos"
	"github.com/stretchr/testify/assert"
)

func TestRoomProxy(t *testing.T) {
	s := &proxy.Statd{}
	// go s.Tick()
	ctx := proxy.WithStatd(context.TODO(), s)
	name := "server1"
	roomproxy, err := NewRoomProxy(ctx, name)
	assert.Nil(t, err)
	target, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:10086")
	for i := 1; i < 32768; i++ {
		roomproxy.Routers().Create(&repos.RouterRecord{Token: uint32(i), Addr: target})
	}
	proxy.NewServer("0.0.0.0:4321", roomproxy).Run(ctx)
}
