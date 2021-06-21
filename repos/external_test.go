package repos

import (
	"fmt"
	"net"
	"testing"

	"github.com/LilithGames/spiracle/config"
	"github.com/stretchr/testify/assert"
)

func TestExternalRepo(t *testing.T) {
	c, err := config.Load("../config.yaml")
	assert.Nil(t, err)
	erepo, err := NewConfigExternalRepo(c)
	assert.Nil(t, err)
	e, err := erepo.Get("dev")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(e.Addrs))
	fmt.Printf("%+v\n", e)
}

func TestUDPAddr(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp4", "www.baidu.com:80")
	assert.Nil(t, err)
	fmt.Printf("%+v\n", addr)
}
