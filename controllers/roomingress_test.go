package controllers

import (
	"testing"
	"time"
	"context"

	v1 "github.com/LilithGames/spiracle/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/LilithGames/spiracle/repos"
	"github.com/stretchr/testify/assert"
)

func TestRoomIngressController(t *testing.T) {
	ctx := context.TODO()
	tokens := repos.NewTsTokenRepo()
	_, err := tokens.Create(ctx, repos.TokenCreationToken(3))
	assert.Nil(t, err)
	rec := &RoomIngressReconciler{TokenRepo: tokens}
	ring1 := &v1.RoomIngress{
		Spec: v1.RoomIngressSpec{
			Rooms: []v1.RoomIngressRoom{
				v1.RoomIngressRoom{
					Id: "room1",
					Server: "server1",
					Upstream: "127.0.0.1:4321",
					Players: []v1.RoomIngressPlayer{
						v1.RoomIngressPlayer{Id: "player1", Token: 0},
						v1.RoomIngressPlayer{Id: "player2", Token: 2},
						v1.RoomIngressPlayer{Id: "player3", Token: 3},
						v1.RoomIngressPlayer{Id: "player4", Token: 0},
					},
				},
			},
		},
	}
	var n int
	var u *time.Duration
	n, u = rec.syncTokens(ring1)
	p1 := GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomId: "room1", PlayerId: "player1"})
	p2 := GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomId: "room1", PlayerId: "player2"})
	p3 := GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomId: "room1", PlayerId: "player3"})
	p4 := GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomId: "room1", PlayerId: "player4"})
	assert.NotNil(t, p1)
	assert.NotNil(t, p2)
	assert.NotEqual(t, int64(0), p1.Player.Token)
	assert.Equal(t, v1.PlayerStatusSuccess, p1.Player.Status)
	assert.Equal(t, int64(2), p2.Player.Token)
	assert.Equal(t, v1.PlayerStatusSuccess, p2.Player.Status)
	assert.Equal(t, int64(3), p3.Player.Token)
	assert.Equal(t, v1.PlayerStatusFailure, p3.Player.Status)
	assert.NotEqual(t, int64(0), p4.Player.Token)
	assert.Equal(t, v1.PlayerStatusSuccess, p4.Player.Status)
	assert.NotEqual(t, 0, n)
	assert.NotNil(t, u)

	n, u = rec.syncTokens(ring1)
	assert.Equal(t, 0, n)
	assert.NotNil(t, u)

	p1 = GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomId: "room1", PlayerId: "player1"})
	p1.Player.Expire = metav1.NewTime(time.Now().UTC().Add(-1*time.Second))
	ring1.Spec.Rooms[0].Players = ring1.Spec.Rooms[0].Players[:3]
	n, u = rec.syncTokens(ring1)
	assert.Equal(t, 2, n)
	p1 = GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomId: "room1", PlayerId: "player1"})
	p4 = GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomId: "room1", PlayerId: "player4"})
	assert.Equal(t, v1.PlayerStatusExpired, p1.Player.Status)
	assert.Nil(t, p4)
}
