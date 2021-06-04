package main

import (
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	ctx := ctrl.SetupSignalHandler()
	db := database(ctx)
	go spiracle(ctx, db)
	go controller(ctx)
	<-ctx.Done()
}
