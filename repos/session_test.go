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
	db, err := db.ProvideServer(ctx, db.ServerLocalConfig())
	assert.Nil(t, err)
	defer db.Shutdown(ctx)
	sr, err := NewSessionRepo("server1", db)
	assert.Nil(t, err)
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:1234")
	err = sr.Create(&Session{Token: 0x01, Src: addr})
	assert.Nil(t, err)
	s, err := sr.Get(0x01)
	assert.Nil(t, err)
	assert.Equal(t, addr, s.Src)
	_, err = sr.Get(0x01)
	assert.True(t, errors.Is(err, ErrKeyNotFound))
}
