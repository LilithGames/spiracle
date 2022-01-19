package main

import (
	"context"
	"log"
	"sync"
	"fmt"
	"time"

	manager "sigs.k8s.io/controller-runtime/pkg/manager"
	// "github.com/buraksezer/olric"
	"github.com/LilithGames/spiracle/proxy"
	"github.com/LilithGames/spiracle/services/roomproxy"
	"github.com/LilithGames/spiracle/repos"
	"github.com/LilithGames/spiracle/config"
)

func spiracle(ctx context.Context, conf *config.Config, wg *sync.WaitGroup, mgr manager.Manager) {
	defer wg.Done()
	s := &proxy.Statd{}
	var th proxy.TickHandler
	if conf.RoomProxy.Debug {
		th = proxy.StdoutTickHandler
	}
	go s.Tick(th)
	// maxproc
	ctx = proxy.WithStatd(ctx, s)
	// sessions, err := repos.NewSessionRepo(db)
	sessions, err := repos.NewSessionRepoV2(repos.SessionMaxIdle(time.Second*time.Duration(conf.RoomProxy.Session.MaxIdleDuration)))
	if err != nil {
		log.Fatalln("create session repo err", err)
	}
	routers := repos.NewK8sRouterRepo(mgr.GetClient())
	for _, s := range conf.RoomProxy.Servers {
		rp, err := roomproxy.NewRoomProxy(ctx, s.Name, roomproxy.RoomProxySessionRepo(sessions), roomproxy.RoomProxyRouterRepo(routers), roomproxy.RoomProxyDebug(conf.RoomProxy.Debug), roomproxy.RoomProxyExpire(time.Duration(conf.RoomProxy.Session.Expire)*time.Second))
		if err != nil {
			log.Fatalln("create roomproxy err", err)
		}
		hostport := fmt.Sprintf("%s:%d", s.Host, s.Port)
		server := proxy.NewServer(hostport, rp)
		for i := 0; i < conf.RoomProxy.Workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := server.Run(ctx); err != nil {
					log.Println("roomproxy server err", err, "name", s.Name, "worker", i)
				}
			}()
		}
	}
}
