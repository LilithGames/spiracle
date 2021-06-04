package repos

import (
	"net"
	"context"
	"testing"

	"github.com/LilithGames/spiracle/infra/db"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	ctx := context.TODO()
	db, err := db.ProvideServer(ctx, db.ServerLocalConfig())
	assert.Nil(t, err)
	defer db.Shutdown(ctx)
	router, err := NewRouterRepo(db)
	assert.Nil(t, err)
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:4321")
	err = router.Create(&RouterRecord{Token: 0x01, Addr: addr, RoomId: "id"})
	assert.Nil(t, err)
	record, err := router.Get(TToken(0x01))
	assert.Nil(t, err)
	assert.Equal(t, "id", record.RoomId)
	count := 0
	err = router.List(func(r *RouterRecord) bool {
		count++
		return true
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestClientRouter(t *testing.T) {
	ctx := context.TODO()
	server, err := db.ProvideServer(ctx, db.ServerLocalConfig())
	assert.Nil(t, err)
	defer server.Shutdown(ctx)

	client, err := db.ProvideClient(ctx, db.ClientLocalConfig())
	assert.Nil(t, err)
	router, err := NewClientRouterRepo(client)
	assert.Nil(t, err)
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:4321")
	err = router.Create(&RouterRecord{Token: 0x01, Addr: addr, RoomId: "id"})
	assert.Nil(t, err)
	record, err := router.Get(TToken(0x01))
	assert.Nil(t, err)
	assert.Equal(t, "id", record.RoomId)
	count := 0
	err = router.List(func(r *RouterRecord) bool {
		count++
		return true
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestMemoryRouter(t *testing.T) {
	router, err := NewMemoryRouterRepo()
	assert.Nil(t, err)
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:4321")
	err = router.Create(&RouterRecord{Token: 0x01, Addr: addr, RoomId: "id"})
	assert.Nil(t, err)
	record, err := router.Get(TToken(0x01))
	assert.Nil(t, err)
	assert.Equal(t, "id", record.RoomId)
	count := 0
	err = router.List(func(r *RouterRecord) bool {
		count++
		return true
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}
