package main

import (
	"context"
	"log"
	"sync"
	"fmt"

	manager "sigs.k8s.io/controller-runtime/pkg/manager"
	"github.com/buraksezer/olric"
	"github.com/LilithGames/spiracle/proxy"
	"github.com/LilithGames/spiracle/services/roomproxy"
	"github.com/LilithGames/spiracle/repos"
	"github.com/LilithGames/spiracle/config"
)

func spiracle(ctx context.Context, conf *config.Config, wg *sync.WaitGroup, db *olric.Olric, mgr manager.Manager) {
	defer wg.Done()
	s := &proxy.Statd{}
	if conf.RoomProxy.Debug {
		go s.Tick()
	}
	// maxproc
	ctx = proxy.WithStatd(ctx, s)
	for _, s := range conf.RoomProxy.Servers {
		rp, err := roomproxy.NewRoomProxy(ctx, s.Name, roomproxy.RoomProxyDb(db), roomproxy.RoomProxyRouterRepo(repos.NewK8sRouterRepo(mgr.GetClient())))
		if err != nil {
			log.Fatalln("create roomproxy err", err)
		}
		hostport := fmt.Sprintf("%s:%d", s.Host, s.Port)
		server := proxy.NewServer(hostport, rp)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := server.Run(ctx); err != nil {
				log.Println("roomproxy server err", err, "name", s.Name)
			}
		}()
	}
}
