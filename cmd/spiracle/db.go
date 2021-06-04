package main

import (
	"context"
	"log"
	"time"

	"github.com/LilithGames/spiracle/infra/db"
	"github.com/buraksezer/olric"
)

func database(ctx context.Context) *olric.Olric {
	conf := db.ServerClusterConfig()
	server, err := db.ProvideServer(ctx, conf)
	if err != nil {
		log.Fatalln("start db err", err)
	}
	go func() {
		<-ctx.Done()
		sctx, _ := context.WithTimeout(context.TODO(), time.Second*10)
		if err := server.Shutdown(sctx); err != nil {
			log.Fatalln("shutdown db err", err)
		}
	}()
	return server
}
