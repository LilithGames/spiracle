package controllers

import (
	"testing"
	v1 "github.com/LilithGames/spiracle/api/v1"

	"github.com/stretchr/testify/assert"
)

func TestDiffRoomStatus(t *testing.T) {
	a := v1.RoomIngressStatus{
		Rooms: []v1.RoomIngressRoomStatus{
			v1.RoomIngressRoomStatus{
				Id: "r1",
				Players: []v1.RoomIngressPlayerStatus{
					v1.RoomIngressPlayerStatus{Id: "p1", Token: 1},
					v1.RoomIngressPlayerStatus{Id: "p2", Token: 2},
				},
			},
			v1.RoomIngressRoomStatus{
				Id: "r2",
				Players: []v1.RoomIngressPlayerStatus{
					v1.RoomIngressPlayerStatus{Id: "p3", Token: 3},
					v1.RoomIngressPlayerStatus{Id: "p4", Token: 4},
				},
			},
		},
	}
	b := v1.RoomIngressStatus{
		Rooms: []v1.RoomIngressRoomStatus{
			v1.RoomIngressRoomStatus{
				Id: "r1",
				Players: []v1.RoomIngressPlayerStatus{
					v1.RoomIngressPlayerStatus{Id: "p1", Token: 7},
					v1.RoomIngressPlayerStatus{Id: "p2", Token: 2},
				},
			},
			v1.RoomIngressRoomStatus{
				Id: "r3",
				Players: []v1.RoomIngressPlayerStatus{
					v1.RoomIngressPlayerStatus{Id: "p5", Token: 5},
					v1.RoomIngressPlayerStatus{Id: "p6", Token: 6},
				},
			},
		},
	}
	r := DiffRoomStatus(a, b)
	assert.Equal(t, 5, len(r))
}
