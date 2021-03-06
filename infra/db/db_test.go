package db

import (
	"context"
	"testing"
	"net"


	"github.com/stretchr/testify/assert"
)

func TestEmbedded(t *testing.T) {
	db, err := ProvideServer(context.TODO(), ServerLocalConfig())
	assert.Nil(t, err)
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:8454")
	items, err := db.NewDMap("items")
	assert.Nil(t, err)
	err = items.Put("name", addr)
	assert.Nil(t, err)
	result, err := items.Get("name")
	assert.Nil(t, err)
	assert.Equal(t, addr, result)
}

func TestClient(t *testing.T) {
	_, err := ProvideServer(context.TODO(), ServerLocalConfig())
	assert.Nil(t, err)
	client, err := ProvideClient(context.TODO(), ClientLocalConfig())
	assert.Nil(t, err)
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:8454")
	items := client.NewDMap("items")
	assert.Nil(t, err)
	err = items.Put("name", addr)
	assert.Nil(t, err)
	result, err := items.Get("name")
	assert.Nil(t, err)
	assert.Equal(t, addr, result)
}

// func serve(c *config.Config, wg *sync.WaitGroup) *olric.Olric {
	// c.Started = func() {
		// defer wg.Done()
		// log.Println("Olric is ready to accept connections")
	// }
	// db, err := olric.New(c)
	// if err != nil {
		// panic(err)
	// }
	// go db.Start()
	// return db
	// defer db.Shutdown(context.TODO())
// }

// func TestOlric(t *testing.T) {
	// wg := &sync.WaitGroup{}
	// c1 := config.New("lan") // lan
	// c1.BindAddr = "127.0.0.1"
	// c1.BindPort = 3300
	// c1.MemberlistConfig.BindAddr = "127.0.0.1"
	// c1.MemberlistConfig.BindPort = 3301
	// c1.Peers = []string{"127.0.0.1:3401"}
	// c2 := config.New("lan") // lan
	// c2.BindAddr = "127.0.0.1"
	// c2.BindPort = 3400
	// c2.MemberlistConfig.BindAddr = "127.0.0.1"
	// c2.MemberlistConfig.BindPort = 3401
	// c2.Peers = []string{"127.0.0.1:3301"}
	// wg.Add(2)
	// db1 := serve(c1, wg)
	// db2 := serve(c2, wg)
	// wg.Wait()

	// addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:8454")
	// items1, err := db1.NewDMap("items")
	// assert.Nil(t, err)
	// err = items1.Put("name", addr)
	// assert.Nil(t, err)

	// items2, err := db2.NewDMap("items")
	// assert.Nil(t, err)
	// v, err := items2.Get("name")
	// assert.Nil(t, err)
	// fmt.Printf("%#v\n", v.(net.Addr))
	// assert.Equal(t, "hulucc", v)
// }

// func BenchmarkDB(b *testing.B) {
	// for i := 0; i < b.N; i++ {
		// _ = i
	// }
// }
