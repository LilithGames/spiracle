package controllers

import (
	"testing"
	"time"
	"context"

	v1 "github.com/LilithGames/spiracle/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/LilithGames/spiracle/repos"
	"github.com/stretchr/testify/assert"
	"github.com/LilithGames/spiracle/config"
)

func TestRoomIngressController(t *testing.T) {
	ctx := context.TODO()
	conf := &config.Config{}
	conf.RoomProxy.Servers = append(conf.RoomProxy.Servers, config.Server{
		Name: "server1",
		Host: "0.0.0.0",
		Port: 5000,
		Externals: []string{"1.2.3.4:7777"},
	})
	externals, err := repos.NewConfigExternalRepo(conf)
	assert.Nil(t, err)

	trepos := make(map[string]repos.TokenRepo)
	trepos["server1"] = repos.NewTsTokenRepo()
	_, err = trepos["server1"].Create(ctx, repos.TokenCreationToken(3))
	assert.Nil(t, err)
	rec := &RoomIngressReconciler{TokenRepos: trepos, ExternalRepos: externals}
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
				v1.RoomIngressRoom{
					Id: "room1",
					Server: "server2",
					Upstream: "127.0.0.1:4321",
					Players: []v1.RoomIngressPlayer{
						v1.RoomIngressPlayer{Id: "player21", Token: 0},
					},
				},
			},
		},
	}
	var n int
	var u *time.Duration
	n, u = rec.syncTokens(ring1)
	room11 := RoomKey{ServerId: "server1", RoomId: "room1"}
	room21 := RoomKey{ServerId: "server2", RoomId: "room1"}
	p1 := GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomKey: room11, PlayerId: "player1"})
	p2 := GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomKey: room11, PlayerId: "player2"})
	p3 := GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomKey: room11, PlayerId: "player3"})
	p4 := GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomKey: room11, PlayerId: "player4"})
	p21 := GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomKey: room21, PlayerId: "player21"})
	assert.NotNil(t, p1)
	assert.NotNil(t, p2)
	assert.NotNil(t, p21)
	assert.NotEqual(t, int64(0), p1.Player.Token)
	assert.Equal(t, v1.PlayerStatusSuccess, p1.Player.Status)
	assert.Equal(t, "1.2.3.4:7777", p1.Player.Externals[0])
	assert.Equal(t, int64(2), p2.Player.Token)
	assert.Equal(t, v1.PlayerStatusSuccess, p2.Player.Status)
	assert.Equal(t, int64(3), p3.Player.Token)
	assert.Equal(t, v1.PlayerStatusFailure, p3.Player.Status)
	assert.NotEqual(t, int64(0), p4.Player.Token)
	assert.Equal(t, v1.PlayerStatusSuccess, p4.Player.Status)
	assert.Equal(t, v1.PlayerStatusFailure, p21.Player.Status)
	assert.NotEqual(t, 0, n)
	assert.NotNil(t, u)

	n, u = rec.syncTokens(ring1)
	assert.Equal(t, 0, n)
	assert.NotNil(t, u)

	p1 = GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomKey: room11, PlayerId: "player1"})
	p1.Player.Expire = ptr(metav1.NewTime(time.Now().UTC().Add(-1*time.Second)))
	ring1.Spec.Rooms[0].Players = ring1.Spec.Rooms[0].Players[:3]
	n, u = rec.syncTokens(ring1)
	assert.Equal(t, 2, n)
	p1 = GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomKey: room11, PlayerId: "player1"})
	p4 = GetPlayerStatusByKey(&ring1.Status, PlayerKey{RoomKey: room11, PlayerId: "player4"})
	assert.Equal(t, v1.PlayerStatusExpired, p1.Player.Status)
	assert.Nil(t, p4)

	ring1.Spec.Rooms[0].Upstream = "127.0.0.2:4321"
	n, u = rec.syncTokens(ring1)
	assert.Equal(t, 3, n)
	// _ = n
	assert.Equal(t, "127.0.0.2:4321", ring1.Status.Rooms[0].Upstream)

	n, _ = rec.syncTokens(ring1)
	assert.Equal(t, 0, n)
}
