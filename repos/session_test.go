package repos

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/LilithGames/spiracle/infra/db"
)

func TestSession(t *testing.T) {
	ctx := context.TODO()
	db, err := db.ProvideLocal(ctx)
	assert.Nil(t, err)
	defer db.Shutdown(ctx)
	sr, err := NewSessionRepo("server1", db)
	assert.Nil(t, err)
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:1234")
	err = sr.Create(&Session{Id: "hulucc", Addr: addr})
	assert.Nil(t, err)
	s, err := sr.Get("hulucc")
	assert.Nil(t, err)
	assert.Equal(t, addr, s.Addr)
	_, err = sr.Get("hulucc1")
	assert.True(t, errors.Is(err, ErrKeyNotFound))
}
