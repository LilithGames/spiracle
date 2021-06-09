package main

import (
	"context"
	"log"
	"time"
	"sync"

	"github.com/LilithGames/spiracle/infra/db"
	"github.com/buraksezer/olric"
)

func database(ctx context.Context, wg *sync.WaitGroup) *olric.Olric {
	conf := db.ServerClusterConfig()
	server, err := db.ProvideServer(ctx, conf)
	if err != nil {
		log.Fatalln("[ERROR] Olric start db err", err)
	}
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
