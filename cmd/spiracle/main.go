package main

import (
	"sync"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
	"github.com/LilithGames/spiracle/config"
)

func main() {
	ctx := ctrl.SetupSignalHandler()
	conf, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalln("load config err: ", err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)
	db := database(ctx, wg)
	mgr := controller(ctx)
	go spiracle(ctx, conf, wg, db, mgr)
	if err := mgr.Start(ctx); err != nil {
		log.Println(err, "[ERROR] controller stop err")
		return
	}
	wg.Wait()
}
