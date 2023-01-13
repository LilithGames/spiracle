package repos

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	v1 "github.com/LilithGames/spiracle/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRouterRepo struct {
	client client.Client
}

var ErrNotSupported = errors.New("notsupported")

func NewK8sRouterRepo(c client.Client) RouterRepo {
	return &k8sRouterRepo{client: c}
}

func (it *k8sRouterRepo) Create(record *RouterRecord, opts ...RouterOption) error {
	return ErrNotSupported
}
func (it *k8sRouterRepo) Update(record *RouterRecord, opts ...RouterOption) error {
	return ErrNotSupported
}
func (it *k8sRouterRepo) CreateOrUpdate(record *RouterRecord, opts ...RouterOption) error {
	return ErrNotSupported
}
func (it *k8sRouterRepo) Delete(token TToken, opts ...RouterOption) error {
	return ErrNotSupported
}
func (it *k8sRouterRepo) Get(token TToken, opts ...RouterOption) (*RouterRecord, error) {
	o := getRouterOptions(opts...)
	indexToken := fmt.Sprintf("%s:%s", o.scope, str(token))
	rings := &v1.RoomIngressList{}
	err := it.client.List(context.TODO(), rings, &client.MatchingFields{"indexToken": indexToken})
	if err != nil {
		return nil, fmt.Errorf("k8s router repo get err: %w", err)
	}
	if len(rings.Items) == 0 {
		return nil, fmt.Errorf("%w: ring not found", ErrNotExists)
	}
	if len(rings.Items) > 1 {
		log.Println("[WARN] duplicated token: ", token)
	}
	ring := &rings.Items[0]
	for i := range ring.Status.Rooms {
		room := &ring.Status.Rooms[i]
		for j := range room.Players {
			player := &room.Players[j]
			if TToken(player.Token) == token {
				addr, err := net.ResolveUDPAddr("udp4", room.Upstream)
				if err != nil {
					return nil, fmt.Errorf("k8s router repo get ResolveUDPAddr err: %w", err)
				}
				return &RouterRecord{
					Token:    token,
					Addr:     addr,
					RoomId:   room.Id,
					PlayerId: player.Id,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("%w: ring: %s", ErrNotExists, ring.ObjectMeta.Name)
}
func (it *k8sRouterRepo) List(f func(*RouterRecord) bool, opts ...RouterOption) error {
	return ErrNotSupported
}

func BuildIndexToken(o client.Object) []string {
	ring := o.(*v1.RoomIngress)
	result := make([]string, 0)
	for i := range ring.Status.Rooms {
		room := &ring.Status.Rooms[i]
		for j := range room.Players {
			player := &room.Players[j]
			if player.Status == v1.PlayerStatusSuccess {
				indexToken := fmt.Sprintf("%s:%s", room.Server, str(uint32(player.Token)))
				result = append(result, indexToken)
			}
		}
	}
	return result
}
