package controllers

import (
	"context"
	"fmt"
	"net"

	v1 "github.com/LilithGames/spiracle/api/v1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/LilithGames/spiracle/repos"
)

const FinalizerName = "projectdavinci.com/finalizer"

type RoomIngressReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Log     logr.Logger
	Routers repos.RouterRepo
}

func (it *RoomIngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := it.Log.WithValues("RoomIngress", req.NamespacedName)
	log.Info("OnReconcile")
	ring := &v1.RoomIngress{}
	if err := it.Get(ctx, req.NamespacedName, ring); err != nil {
		if client.IgnoreNotFound(err) == nil {
			log.Info("RoomIngress deleted")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("get err: %w", err)
	}
	if ring.ObjectMeta.DeletionTimestamp.IsZero() {
		if !contains(ring.GetFinalizers(), FinalizerName) {
			controllerutil.AddFinalizer(ring, FinalizerName)
			it.syncTokens(ring)
			if err := it.Update(ctx, ring); err != nil {
				return ctrl.Result{}, fmt.Errorf("AddFinalizer update err: %w", err)
			}
		}
	} else if contains(ring.GetFinalizers(), FinalizerName) {
		log.Info("removing external resource")
		controllerutil.RemoveFinalizer(ring, FinalizerName)
		ring.Spec = v1.RoomIngressSpec{}
		it.syncTokens(ring)
		if err := it.Update(ctx, ring); err != nil {
			return ctrl.Result{}, fmt.Errorf("RemoveFinalizer update err: %w", err)
		}
	}
	if n := it.syncTokens(ring); n > 0 {
		if err := it.Update(ctx, ring); err != nil {
			return ctrl.Result{}, fmt.Errorf("update err: %w", err)
		}
	}
	return ctrl.Result{}, nil
}

func (it *RoomIngressReconciler) getRouterRecord(s *v1.RoomIngressStatus, pos PlayerPos) (*repos.RouterRecord, string) {
	room := &s.Rooms[pos.RoomIndex]
	player := &room.Players[pos.PlayerIndex]
	addr, _ := net.ResolveUDPAddr("udp4", room.Upstream)
	return &repos.RouterRecord{
		Token:    uint32(player.Token),
		Addr:     addr,
		RoomId:   room.Id,
		PlayerId: player.Id,
	}, room.Server
}

func (it *RoomIngressReconciler) syncTokens(ring *v1.RoomIngress) int {
	past := ring.Status
	curr := GetStatus(ring.Spec)
	diffs := DiffRoomStatus(past, curr)
	for i := range diffs {
		diff := &diffs[i]
		switch diff.Type {
		case DiffNew:
			record, scope := it.getRouterRecord(&curr, diff.Current)
			if err := it.Routers.Create(record, repos.RouterScope(scope)); err != nil {
				it.Log.Error(err, "create router err", "record", record, "scope", scope)
			}
		case DiffUpdated:
			crecord, cscope := it.getRouterRecord(&curr, diff.Current)
			precord, pscope := it.getRouterRecord(&past, diff.Past)
			if err := it.Routers.Create(crecord, repos.RouterScope(cscope)); err != nil {
				it.Log.Error(err, "create router err", "record", crecord, "scope", cscope)
			}
			if err := it.Routers.Delete(precord.Token, repos.RouterScope(pscope)); err != nil {
				it.Log.Error(err, "delete router err", "record", precord, "scope", pscope)
			}
		case DiffDeleted:
			record, scope := it.getRouterRecord(&past, diff.Past)
			if err := it.Routers.Delete(record.Token, repos.RouterScope(scope)); err != nil {
				it.Log.Error(err, "delete router err", "record", record, "scope", scope)
			}
		default:
			panic("unknown diff type")
		}
	}
	ring.Status = curr
	return len(diffs)
}

func (it *RoomIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	pred := predicate.GenerationChangedPredicate{}
	return ctrl.NewControllerManagedBy(mgr).For(&v1.RoomIngress{}).WithEventFilter(pred).Complete(it)
}
