package main

import (
	"flag"
	"sync"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
	"github.com/LilithGames/spiracle/config"
)

func main() {
	cpath := flag.String("config", "config.yaml", "config path")
	flag.Parse()
	ctx := ctrl.SetupSignalHandler()
	wg := &sync.WaitGroup{}
	conf, err := config.Load(*cpath)
	if err != nil {
		log.Fatalln("load config err: ", err)
	}

	mgr := controller(ctx, conf)

	if conf.RoomProxy.Enable {
		db := database(ctx, wg, conf)
		wg.Add(1)
		go spiracle(ctx, conf, wg, db, mgr)
	}

	if err := mgr.Start(ctx); err != nil {
		log.Println(err, "[ERROR] controller stop err")
		return
	}
	wg.Wait()
}
