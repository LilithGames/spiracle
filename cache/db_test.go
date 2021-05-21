package cache

import (
	"log"
	"sync"
	"testing"

	// "time"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"
	"github.com/stretchr/testify/assert"
)

func serve(c *config.Config, wg *sync.WaitGroup) *olric.Olric {
	c.Started = func() {
		defer wg.Done()
		log.Println("Olric is ready to accept connections")
	}
	db, err := olric.New(c)
	if err != nil {
		panic(err)
	}
	go db.Start()
	return db
	// defer db.Shutdown(context.TODO())
}

func TestOlric(t *testing.T) {
	wg := &sync.WaitGroup{}
	c1 := config.New("lan") // lan
	c1.BindAddr = "127.0.0.1"
	c1.BindPort = 3300
	c1.MemberlistConfig.BindAddr = "127.0.0.1"
	c1.MemberlistConfig.BindPort = 3301
	c1.Peers = []string{"127.0.0.1:3401"}
	c2 := config.New("lan") // lan
	c2.BindAddr = "127.0.0.1"
	c2.BindPort = 3400
	c2.MemberlistConfig.BindAddr = "127.0.0.1"
	c2.MemberlistConfig.BindPort = 3401
	c2.Peers = []string{"127.0.0.1:3301"}
	wg.Add(2)
	db1 := serve(c1, wg)
	db2 := serve(c2, wg)
	wg.Wait()

	items1, err := db1.NewDMap("items")
	assert.Nil(t, err)
	err = items1.Put("name", "hulucc")

	items2, err := db2.NewDMap("items")
	assert.Nil(t, err)
	v, err := items2.Get("name")
	assert.Nil(t, err)
	assert.Equal(t, "hulucc", v)
}

func BenchmarkDB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = i
	}
}
