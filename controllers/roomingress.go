package controllers

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/LilithGames/spiracle/api/v1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/LilithGames/spiracle/repos"
	"github.com/LilithGames/spiracle/config"
)

const FinalizerName = "projectdavinci.com/finalizer"

type RoomIngressReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Log           logr.Logger
	TokenRepos    map[string]repos.TokenRepo
	ExternalRepos repos.ExternelRepo
	Config        *config.Config
}

func (it *RoomIngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := it.Log.WithValues("RoomIngress", req.NamespacedName)
	start := time.Now()
	log.Info("OnReconcile begin")
	r, err := it.reconcile(ctx, req)
	log.Info("OnReconcile end", "duration", time.Since(start))
	return r, err
}

func (it *RoomIngressReconciler) reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := it.Log.WithValues("RoomIngress", req.NamespacedName)
	ring := &v1.RoomIngress{}
	var curr *v1.RoomIngress
	patch := func() *v1.RoomIngress {
		if curr == nil {
			curr = ring.DeepCopy()
		}
		return curr
	}
	if err := it.Get(ctx, req.NamespacedName, ring); err != nil {
		if client.IgnoreNotFound(err) == nil {
			log.Info("RoomIngress deleted")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("get err: %w", err)
	}
	if ring.ObjectMeta.DeletionTimestamp.IsZero() && it.Config.Controller.Reconciler.Finalizer.Enable {
		if !contains(ring.GetFinalizers(), FinalizerName) {
			controllerutil.AddFinalizer(patch(), FinalizerName)
			if err := it.Patch(ctx, patch(), client.MergeFrom(ring)); err != nil {
				return ctrl.Result{}, fmt.Errorf("AddFinalizer update err: %w", err)
			}
		}
	} else if contains(ring.GetFinalizers(), FinalizerName) {
		log.Info("removing external resource")
		controllerutil.RemoveFinalizer(patch(), FinalizerName)
		patch().Spec = v1.RoomIngressSpec{}
		if err := it.Patch(ctx, patch(), client.MergeFrom(ring)); err != nil {
			return ctrl.Result{}, fmt.Errorf("RemoveFinalizer update err: %w", err)
		}
		return ctrl.Result{}, nil
	}
	if n, requeue := it.syncTokens(patch()); n > 0 {
		if err := it.Status().Patch(ctx, patch(), client.MergeFrom(ring)); err != nil {
			return it.requeue(requeue), fmt.Errorf("update status err: %w", err)
		}
		return it.requeue(requeue), nil
	}
	return ctrl.Result{}, nil
}

func (it *RoomIngressReconciler) requeue(d *time.Duration) ctrl.Result {
	if d != nil {
		return ctrl.Result{Requeue: true, RequeueAfter: *d}
	}
	return ctrl.Result{}
}

func (it *RoomIngressReconciler) syncTokens(ring *v1.RoomIngress) (int, *time.Duration) {
	requeue := make([]time.Duration, 0)
	past := &ring.Status
	curr := GetStatus(&ring.Spec, past)
	diffs := DiffRoomStatus(past, curr, DiffUpdater(TokenUpdatedHandler()))
	n := 0
	for i := range diffs {
		diff := &diffs[i]
		if diff.Type == DiffNew || diff.Type == DiffUpdated {
			c := GetPlayerStatusByPos(curr, diff.Current)
			if c.Player.Status == v1.PlayerStatusExpired {
				continue
			}
			repo, ok := it.TokenRepos[c.Room.Server]
			if !ok {
				c.Player.Status = v1.PlayerStatusFailure
				c.Player.Detail = "unknown room.server"
				n++
				continue
			}
			external, err := it.ExternalRepos.Get(c.Room.Server)
			if err != nil {
				c.Player.Status = v1.PlayerStatusFailure
				c.Player.Detail = "get server external err: " + err.Error()
				n++
				continue
			}
			token, err := repo.Create(context.TODO(), repos.TokenCreationToken(uint32(c.Player.Token)))
			if err != nil {
				c.Player.Status = v1.PlayerStatusFailure
				c.Player.Detail = err.Error()
				n++
				continue
				// requeue = append(requeue, time.Minute)
			}
			c.Player.Status = v1.PlayerStatusSuccess
			c.Player.Token = int64(token.TToken)
			c.Player.Externals = external.HostPorts()
			c.Player.Timestamp = ptr(metav1.NewTime(token.Timestamp))
			c.Player.Expire = ptr(metav1.NewTime(token.Expire))
			n++
			requeue = append(requeue, token.Duration())
		}
		if diff.Type == DiffDeleted || diff.Type == DiffUpdated {
			p := GetPlayerStatusByPos(past, diff.Past)
			if p.Player.Status == v1.PlayerStatusSuccess {
				if repo, ok := it.TokenRepos[p.Room.Server]; ok {
					repo.Delete(context.TODO(), uint32(p.Player.Token))
				}
			}
			if diff.Type == DiffDeleted {
				n++
			}
		}
		if diff.Type == DiffUnchanged {
			c := GetPlayerStatusByPos(curr, diff.Current)
			if c.Player.Status == v1.PlayerStatusSuccess {
				now := time.Now().UTC()
				expire := c.Player.Expire.Time
				if now.After(expire) {
					c.Player.Status = v1.PlayerStatusExpired
					c.Player.Detail = "token expired"
					if repo, ok := it.TokenRepos[c.Room.Server]; ok {
						repo.Delete(context.TODO(), uint32(c.Player.Token))
					}
					n++
				} else {
					requeue = append(requeue, expire.Sub(now))
				}
			}
		}
	}
	ring.Status = *curr
	return n, min(requeue)
}

func (it *RoomIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	pred := predicate.GenerationChangedPredicate{}
	opts := controller.Options{
		MaxConcurrentReconciles: it.Config.Controller.Reconciler.Concurrency,
	}
	return ctrl.NewControllerManagedBy(mgr).For(&v1.RoomIngress{}).WithEventFilter(pred).WithOptions(opts).Complete(it)
}
