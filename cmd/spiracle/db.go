package main

import (
	"context"
	"log"
	"time"
	"sync"

	"github.com/LilithGames/spiracle/infra/db"
	"github.com/LilithGames/spiracle/config"
	"github.com/buraksezer/olric"
)

func database(ctx context.Context, wg *sync.WaitGroup, conf *config.Config) *olric.Olric {
	dbconf := db.ServerClusterConfig(conf.RoomProxy.Session.MaxIdleDuration)
	server, err := db.ProvideServer(ctx, dbconf)
	if err != nil {
		log.Fatalln("[ERROR] Olric start db err", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		sctx, _ := context.WithTimeout(context.TODO(), time.Second*10)
		if err := server.Shutdown(sctx); err != nil {
			log.Println("[ERROR] Olric shutdown err", err)
			return
		}
		log.Println("[INFO] Olric is shutdown gracefully")
	}()
	return server
}
