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
	db, err := db.ProvideLocal(ctx)
	assert.Nil(t, err)
	defer db.Shutdown(ctx)
	router, err := NewRouterRepo("server1", db)
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
