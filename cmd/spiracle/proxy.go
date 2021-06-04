package main

import (
	"context"
	"log"

	"github.com/buraksezer/olric"
	"github.com/LilithGames/spiracle/proxy"
	"github.com/LilithGames/spiracle/services/roomproxy"
)

func spiracle(ctx context.Context, db *olric.Olric) {
	s := &proxy.Statd{}
	// go s.Tick()
	// maxproc
	ctx = proxy.WithStatd(ctx, s)
	rp, err := roomproxy.NewRoomProxy(ctx, "roomproxy", roomproxy.RoomProxyDb(db))
	if err != nil {
		log.Fatalln("create roomproxy err", err)
	}
	server := proxy.NewServer("0.0.0.0:4321", rp)
	if err := server.Run(ctx); err != nil {
		log.Println("roomproxy server err", err)
	}
}
