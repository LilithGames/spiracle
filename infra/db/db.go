package db

import (
	"context"
	"errors"
	"log"
	"time"
	"os"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"
	"github.com/buraksezer/olric/client"
	discovery "github.com/buraksezer/olric-cloud-plugin/lib"
)

func ServerLocalConfig() *config.Config {
	c := config.New("local")
	c.BindAddr = "127.0.0.1"
	c.MemberlistConfig.BindAddr = "127.0.0.1"
	return c
}

func ServerClusterConfig() *config.Config {
	c := config.New("lan")
	if os.Getenv("OLRIC_DISCOVERY_PROVIDER") == "k8s" {
		ns := os.Getenv("OLRIC_DISCOVERY_NAMESPACE")
		labelname := os.Getenv("OLRIC_DISCOVERY_LABEL_NAME")
		labelvalue := os.Getenv("OLRIC_DISCOVERY_LABEL_VALUE")
		labelSelector := fmt.Sprintf("%s=%s", labelname, labelvalue)
		c.ServiceDiscovery = map[string]interface{}{
			"plugin": &discovery.CloudDiscovery{},
			"provider": "k8s",
			"args": fmt.Sprintf("namespace=%s label_selector=\"%s\"", ns, labelSelector)
		}
	}
	return c
}

func ClientLocalConfig() *client.Config {
	c := &client.Config{
		Servers:       []string{"localhost:3320"},
		Client: &config.Client{
			DialTimeout: 10 * time.Second,
			KeepAlive:   10 * time.Second,
			MaxConn:     100,
		},
	}
	return c
}

func ProvideServer(ctx context.Context, c *config.Config) (*olric.Olric, error) {
	ready := make(chan struct{})
	done := make(chan error)
	c.Started = func() {
		defer close(ready)
		log.Println("[INFO] Olric is ready to accept connections")
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

func ProvideClient(ctx context.Context, c *client.Config) (*client.Client, error) {
	client, err := client.New(c)
	if err != nil {
		return nil, err
	}
	return client, nil
}
