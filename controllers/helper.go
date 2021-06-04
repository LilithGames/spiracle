package controllers

import v1 "github.com/LilithGames/spiracle/api/v1"

func GetPlayerStatus(player v1.RoomIngressPlayer) v1.RoomIngressPlayerStatus {
	return v1.RoomIngressPlayerStatus{
		Id:    player.Id,
		Token: player.Token,
	}
}

func GetRoomStatus(room v1.RoomIngressRoom) v1.RoomIngressRoomStatus {
	players := make([]v1.RoomIngressPlayerStatus, len(room.Players))
	for i := range room.Players {
		players[i] = GetPlayerStatus(room.Players[i])
	}
	return v1.RoomIngressRoomStatus{
		Id:       room.Id,
		Server:   room.Server,
		Upstream: room.Upstream,
		Players:  players,
	}
}

func GetStatus(spec v1.RoomIngressSpec) v1.RoomIngressStatus {
	rooms := make([]v1.RoomIngressRoomStatus, len(spec.Rooms))
	for i := range spec.Rooms {
		rooms[i] = GetRoomStatus(spec.Rooms[i])
	}
	return v1.RoomIngressStatus{Rooms: rooms}
}

type DiffType string

const DiffNew = "New"
const DiffUpdated = "Updated"
const DiffDeleted = "Deleted"

type PlayerKey struct {
	RoomId   string
	PlayerId string
}

type PlayerPos struct {
	RoomIndex   int
	PlayerIndex int
}

var NilPos = PlayerPos{RoomIndex: -1, PlayerIndex: -1}

type DiffResult struct {
	Type    DiffType
	Key     PlayerKey
	Past    PlayerPos
	Current PlayerPos
}

func GetPlayerStatusDict(s v1.RoomIngressStatus) map[PlayerKey]PlayerPos {
	r := make(map[PlayerKey]PlayerPos)
	for i := range s.Rooms {
		room := &s.Rooms[i]
		for j := range room.Players {
			player := &room.Players[j]
			key := PlayerKey{RoomId: room.Id, PlayerId: player.Id}
			r[key] = PlayerPos{RoomIndex: i, PlayerIndex: j}
		}
	}
	return r
}

func GetPlayerStatusDictKeys(dict map[PlayerKey]PlayerPos) []PlayerKey {
	r := make([]PlayerKey, 0, len(dict))
	for k, _ := range dict {
		r = append(r, k)
	}
	return r
}

func UnionPlayerKeys(a map[PlayerKey]PlayerPos, b map[PlayerKey]PlayerPos) []PlayerKey {
	dict := make(map[PlayerKey]struct{})
	for k, _ := range a {
		dict[k] = struct{}{}
	}
	for k, _ := range b {
		dict[k] = struct{}{}
	}
	r := make([]PlayerKey, 0, len(dict))
	for k, _ := range dict {
		r = append(r, k)
	}
	return r
}

func GetPlayerStatusByPos(s v1.RoomIngressStatus, pos PlayerPos) v1.RoomIngressPlayerStatus {
	return s.Rooms[pos.RoomIndex].Players[pos.PlayerIndex]
}

func DiffRoomStatus(past v1.RoomIngressStatus, curr v1.RoomIngressStatus, opts ...DiffOption) []DiffResult {
	o := getDiffOptions(opts...)
	pastd := GetPlayerStatusDict(past)
	currd := GetPlayerStatusDict(curr)
	keys := UnionPlayerKeys(pastd, currd)
	r := make([]DiffResult, 0, len(keys))
	for _, key := range keys {
		diff := DiffResult{
			Key:     key,
			Current: NilPos,
			Past:    NilPos,
		}
		p, pok := pastd[key]
		c, cok := currd[key]
		if cok && !pok {
			diff.Type = DiffNew
			diff.Current = c
		} else if cok && pok {
			if o.uh(past, p, curr, c) {
				diff.Type = DiffUpdated
				diff.Current = c
				diff.Past = p
			} else {
				continue
			}
		} else if !cok && pok {
			diff.Type = DiffDeleted
			diff.Past = p
		} else {
			panic("impossible")
		}
		r = append(r, diff)
	}
	return r
}

func TokenUpdatedHandler() UpdatedHandler {
	return func(past v1.RoomIngressStatus, pp PlayerPos, curr v1.RoomIngressStatus, cp PlayerPos) bool {
		pps := GetPlayerStatusByPos(past, pp)
		cps := GetPlayerStatusByPos(curr, cp)
		return pps.Token != cps.Token
	}
}

type UpdatedHandler func(past v1.RoomIngressStatus, pp PlayerPos, curr v1.RoomIngressStatus, cp PlayerPos) bool

type diffOptions struct {
	uh UpdatedHandler
}

type DiffOption interface {
	apply(*diffOptions)
}

type funcDiffOption struct {
	f func(*diffOptions)
}

func (it *funcDiffOption) apply(o *diffOptions) {
	it.f(o)
}

func newFuncDiffOption(f func(*diffOptions)) DiffOption {
	return &funcDiffOption{f: f}
}
func getDiffOptions(opts ...DiffOption) *diffOptions {
	o := &diffOptions{
		uh: TokenUpdatedHandler(),
	}
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}
