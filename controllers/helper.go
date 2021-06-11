package controllers

import (
	v1 "github.com/LilithGames/spiracle/api/v1"
)

func GetPlayerStatus(player v1.RoomIngressPlayer, base *v1.RoomIngressPlayerStatus) v1.RoomIngressPlayerStatus {
	s := v1.RoomIngressPlayerStatus{
		Id: player.Id,
	}
	if base == nil {
		s.Token = player.Token
		s.Status = v1.PlayerStatusPending
	} else {
		if player.Token != int64(0) {
			s.Token = player.Token
		} else {
			s.Token = base.Token
		}
		s.Status = base.Status
		s.Detail = base.Detail
		s.Timestamp = base.Timestamp
		s.Expire = base.Expire
	}
	return s
}

func GetRoomStatus(room v1.RoomIngressRoom, base *v1.RoomIngressRoomStatus) v1.RoomIngressRoomStatus {
	players := make([]v1.RoomIngressPlayerStatus, len(room.Players))
	dict := make(map[string]*v1.RoomIngressPlayerStatus)
	if base != nil {
		for i := range base.Players {
			player := &base.Players[i]
			dict[player.Id] = player
		}
	}
	for i := range room.Players {
		pbase, _ := dict[room.Players[i].Id]
		players[i] = GetPlayerStatus(room.Players[i], pbase)
	}
	return v1.RoomIngressRoomStatus{
		Id:       room.Id,
		Server:   room.Server,
		Upstream: room.Upstream,
		Players:  players,
	}
}

func GetRoomStatusDict(rooms []v1.RoomIngressRoomStatus) map[RoomKey]*v1.RoomIngressRoomStatus {
	dict := make(map[RoomKey]*v1.RoomIngressRoomStatus)
	for i := range rooms {
		room := &rooms[i]
		dict[RoomKey{ServerId: room.Server, RoomId: room.Id}] = room
	}
	return dict
}

func GetStatus(spec *v1.RoomIngressSpec, base *v1.RoomIngressStatus) *v1.RoomIngressStatus {
	rooms := make([]v1.RoomIngressRoomStatus, len(spec.Rooms))
	var dict map[RoomKey]*v1.RoomIngressRoomStatus
	if base != nil {
		dict = GetRoomStatusDict(base.Rooms)
	}
	for i := range spec.Rooms {
		room := &spec.Rooms[i]
		key := RoomKey{ServerId: room.Server, RoomId: room.Id}
		rbase := dict[key]
		rooms[i] = GetRoomStatus(spec.Rooms[i], rbase)
	}
	return &v1.RoomIngressStatus{Rooms: rooms}
}

type DiffType string

const DiffNew = "New"
const DiffUpdated = "Updated"
const DiffDeleted = "Deleted"
const DiffUnchanged = "Unchanged"

type RoomKey struct {
	ServerId string
	RoomId   string
}
type PlayerKey struct {
	RoomKey
	PlayerId string
}

type PlayerPos struct {
	RoomIndex   int
	PlayerIndex int
}

type PlayerDetail struct {
	Room   *v1.RoomIngressRoomStatus
	Player *v1.RoomIngressPlayerStatus
}

var NilPos = PlayerPos{RoomIndex: -1, PlayerIndex: -1}

type DiffResult struct {
	Type    DiffType
	Key     PlayerKey
	Past    PlayerPos
	Current PlayerPos
}

func GetPlayerStatusDict(s *v1.RoomIngressStatus) map[PlayerKey]PlayerPos {
	r := make(map[PlayerKey]PlayerPos)
	for i := range s.Rooms {
		room := &s.Rooms[i]
		roomkey := RoomKey{ServerId: room.Server, RoomId: room.Id}
		for j := range room.Players {
			player := &room.Players[j]
			key := PlayerKey{RoomKey: roomkey, PlayerId: player.Id}
			r[key] = PlayerPos{RoomIndex: i, PlayerIndex: j}
		}
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

func GetPlayerStatusByPos(s *v1.RoomIngressStatus, pos PlayerPos) PlayerDetail {
	room := &s.Rooms[pos.RoomIndex]
	player := &room.Players[pos.PlayerIndex]
	return PlayerDetail{Room: room, Player: player}
}

func GetPlayerStatusByKey(s *v1.RoomIngressStatus, key PlayerKey) *PlayerDetail {
	for i := range s.Rooms {
		room := &s.Rooms[i]
		if room.Server == key.ServerId && room.Id == key.RoomId {
			for j := range room.Players {
				player := &room.Players[j]
				if player.Id == key.PlayerId {
					return &PlayerDetail{Room: room, Player: player}
				}
			}
		}
	}
	return nil
}

func DiffRoomStatus(past *v1.RoomIngressStatus, curr *v1.RoomIngressStatus, opts ...DiffOption) []DiffResult {
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
			} else {
				diff.Type = DiffUnchanged
			}
			diff.Current = c
			diff.Past = p
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
	return func(past *v1.RoomIngressStatus, pp PlayerPos, curr *v1.RoomIngressStatus, cp PlayerPos) bool {
		p := GetPlayerStatusByPos(past, pp)
		c := GetPlayerStatusByPos(curr, cp)
		if c.Player.Status == v1.PlayerStatusRetry {
			return true
		}
		return p.Player.Token != c.Player.Token
	}
}

func AlwaysUpdatedHandler() UpdatedHandler {
	return func(past *v1.RoomIngressStatus, pp PlayerPos, curr *v1.RoomIngressStatus, cp PlayerPos) bool {
		return true
	}
}

type UpdatedHandler func(past *v1.RoomIngressStatus, pp PlayerPos, curr *v1.RoomIngressStatus, cp PlayerPos) bool

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

func DiffUpdater(uh UpdatedHandler) DiffOption {
	return newFuncDiffOption(func(o *diffOptions) {
		o.uh = uh
	})
}
