package db

import (
	"context"
	"errors"
	"log"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"
)

func ProvideLocal(ctx context.Context) (*olric.Olric, error) {
	ready := make(chan struct{})
	done := make(chan error)
	c := config.New("local")
	c.BindAddr = "127.0.0.1"
	c.MemberlistConfig.BindAddr = "127.0.0.1"
	c.Started = func() {
		defer close(ready)
		log.Println("Olric is ready to accept connections")
	}
	db, err := olric.New(c)
	if err != nil {
		return nil, err
	}
	go func() {
		defer close(done)
		done <- db.Start()
	}()
	select {
	case <-ready:
		return db, nil
	case err, ok := <-done:
		if !ok {
			return nil, errors.New("unknown error")
		}
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
