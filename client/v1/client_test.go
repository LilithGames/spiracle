package v1

import (
	"testing"
	"context"
	"fmt"
	"sync"
	"time"

	spv1 "github.com/LilithGames/spiracle/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var mgr manager.Manager
func init() {
	scheme := runtime.NewScheme()
	utilruntime.Must(spv1.AddToScheme(scheme))
	var err error
	config := ctrl.GetConfigOrDie()
	config.QPS = 1000
	config.Burst = 2000
	mgr, err = ctrl.NewManager(config, ctrl.Options{Scheme: scheme, MetricsBindAddress: ":0", HealthProbeBindAddress: ":0"})
	if err != nil {
		panic(err)
	}
	go mgr.Start(context.TODO())
}

func TestClient(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		func(room string) {
			defer wg.Done()
			createRoom(t, room)
		}(fmt.Sprintf("room%d", i))
	}
	wg.Wait()
}

func TestClientPatch(t *testing.T) {
	wg := sync.WaitGroup{}
	worker := func(id int) {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			name := fmt.Sprintf("room%d", id*10+i%10+1)
			room := spv1.RoomIngressRoom{
				Id: name,
				Server: "local",
				Upstream: "127.0.0.1:9200",
				Players: []spv1.RoomIngressPlayer{
					spv1.RoomIngressPlayer{Id: "Player", Token: 0},
				},
			}
			ring := &spv1.RoomIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
					Namespace: "default",
				},
				Spec: spv1.RoomIngressSpec{
					Rooms: []spv1.RoomIngressRoom{room},
				},
			}
			start := time.Now()
			if err := mgr.GetClient().Status().Patch(context.TODO(), ring, client.MergeFrom(ring)); err != nil {
				fmt.Printf("%+v\n", err)
			}
			fmt.Printf("patch %s duration: %+v\n", name, time.Since(start))
		}
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(i)
	}
	wg.Wait()
}


func createRoom(t *testing.T, name string) {
	room := spv1.RoomIngressRoom{
		Id: name,
		Server: "local",
		Upstream: "127.0.0.1:9200",
		Players: []spv1.RoomIngressPlayer{},
	}
	ring := &spv1.RoomIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Namespace: "default",
		},
		Spec: spv1.RoomIngressSpec{
			Rooms: []spv1.RoomIngressRoom{room},
		},
	}
	err := mgr.GetClient().Create(context.TODO(), ring)
	assert.Nil(t, err)
}

func deleteRoom(room string) {

}

func createPlayer(room string, player string) {

}

func deletePlayer(room string, player string) {

}
