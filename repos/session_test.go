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
	sr, err := NewSessionRepo(db)
	assert.Nil(t, err)
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:1234")
	err = sr.Create(&Session{Token: 0x01, Src: addr}, SessionScope("s1"))
	assert.Nil(t, err)
	s, err := sr.Get(0x01, SessionScope("s1"))
	assert.Nil(t, err)
	assert.Equal(t, addr, s.Src)
	_, err = sr.Get(0x02, SessionScope("s1"))
	assert.True(t, errors.Is(err, ErrNotExists))
}

func TestSessionMemory(t *testing.T) {
	sr := NewMemorySessionRepo()
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:1234")
	err := sr.Create(&Session{Token: 0x01, Src: addr}, SessionScope("s1"))
	assert.Nil(t, err)
	s, err := sr.Get(0x01, SessionScope("s1"))
	assert.Nil(t, err)
	assert.Equal(t, addr, s.Src)
	_, err = sr.Get(0x02, SessionScope("s1"))
	assert.True(t, errors.Is(err, ErrNotExists))
}

func TestSessionV2(t *testing.T) {
	sr, err := NewSessionRepoV2()
	assert.Nil(t, err)
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:1234")
	err = sr.CreateOrUpdate(&Session{Token: 0x01, Src: addr}, SessionScope("s1"))
	assert.Nil(t, err)
	s, err := sr.Get(0x01, SessionScope("s1"))
	assert.Nil(t, err)
	assert.Equal(t, addr, s.Src)
	_, err = sr.Get(0x02, SessionScope("s1"))
	assert.True(t, errors.Is(err, ErrNotExists))
	_, err = sr.Get(0x01, SessionScope("s2"))
	assert.True(t, errors.Is(err, ErrNotExists))
}
