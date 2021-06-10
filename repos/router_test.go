package repos

import (
	"net"
	"context"
	"testing"
	"time"
	"fmt"

	"github.com/LilithGames/spiracle/infra/db"
	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	v1 "github.com/LilithGames/spiracle/api/v1"
)

func TestRouter(t *testing.T) {
	ctx := context.TODO()
	db, err := db.ProvideServer(ctx, db.ServerLocalConfig())
	assert.Nil(t, err)
	defer db.Shutdown(ctx)
	router, err := NewRouterRepo(db)
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

func TestClientRouter(t *testing.T) {
	ctx := context.TODO()
	server, err := db.ProvideServer(ctx, db.ServerLocalConfig())
	assert.Nil(t, err)
	defer server.Shutdown(ctx)

	client, err := db.ProvideClient(ctx, db.ClientLocalConfig())
	assert.Nil(t, err)
	router := NewClientRouterRepo(client)
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

func TestMemoryRouter(t *testing.T) {
	router := NewMemoryRouterRepo()
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:4321")
	err = router.Create(&RouterRecord{Token: 0x01, Addr: addr, RoomId: "id"}, RouterScope("test1"))
	assert.Nil(t, err)
	record, err := router.Get(TToken(0x01), RouterScope("test1"))
	assert.Nil(t, err)
	assert.Equal(t, "id", record.RoomId)
	count := 0
	err = router.List(func(r *RouterRecord) bool {
		count++
		return true
	}, RouterScope("test1"))
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestK8sRouter(t *testing.T) {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{Scheme: scheme, MetricsBindAddress: ":0", HealthProbeBindAddress: ":0"})
	assert.Nil(t, err)
	ctx := context.TODO()
	mgr.GetFieldIndexer().IndexField(ctx, &v1.RoomIngress{}, "indexToken", BuildIndexToken)
	go mgr.Start(ctx)
	time.Sleep(time.Second)
	repo := NewK8sRouterRepo(mgr.GetClient())
	record, err := repo.Get(uint32(4), RouterScope("local"))
	assert.Nil(t, err)
	fmt.Printf("%+v\n", record)
}
