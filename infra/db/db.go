package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	// "strconv"

	"github.com/buraksezer/olric"
	discovery "github.com/buraksezer/olric-cloud-plugin/lib"
	"github.com/buraksezer/olric/client"
	"github.com/buraksezer/olric/config"
)

func ServerLocalConfig() *config.Config {
	c := config.New("local")
	c.BindAddr = "127.0.0.1"
	c.MemberlistConfig.BindAddr = "127.0.0.1"
	return c
}

func ServerClusterConfig() *config.Config {
	c := config.New("lan")
	c.BindAddr = "0.0.0.0"
	c.MemberlistConfig.BindAddr = "0.0.0.0"
	c.DMaps.Custom["cache.sessions"] = config.DMap{MaxIdleDuration: time.Second*30}
	// c.ReplicationMode = config.AsyncReplicationMode
	// c.ReplicaCount = 1
	// c.ReadRepair = true
	// c.MaxJoinAttempts = 10
	// c.WriteQuorum = 1
	// c.ReadQuorum = 1
	// c.MemberCountQuorum = 2
	if os.Getenv("OLRIC_DISCOVERY_PROVIDER") == "k8s" {
		ns := os.Getenv("OLRIC_DISCOVERY_NAMESPACE")
		labelname := os.Getenv("OLRIC_DISCOVERY_LABEL_NAME")
		labelvalue := os.Getenv("OLRIC_DISCOVERY_LABEL_VALUE")
		labelSelector := fmt.Sprintf("%s=%s", labelname, labelvalue)
		c.ServiceDiscovery = map[string]interface{}{
			"plugin":   &discovery.CloudDiscovery{},
			"provider": "k8s",
			"args":     fmt.Sprintf("namespace=%s label_selector=\"%s\"", ns, labelSelector),
		}
	}
	return c
}

func ClientLocalConfig() *client.Config {
	c := &client.Config{
		Servers: []string{"localhost:3320"},
		Client: &config.Client{
			DialTimeout: 10 * time.Second,
			KeepAlive:   10 * time.Second,
			MaxConn:     100,
		},
	}
	return c
}

func ClientClusterLocalConfig() *client.Config {
	host := os.Getenv("OLRIC_CLIENT_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("OLRIC_CLIENT_PORT")
	if port == "" {
		port = "3320"
	}
	c := &client.Config{
		Servers: []string{fmt.Sprintf("%s:%s", host, port)},
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
		log.Println("[INFO] Olric is started")
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
		dmap, err := db.NewDMap("status.ready")
		if err != nil {
			return nil, err
		}
		hostname, _ := os.Hostname()
		if err := dmap.Put(hostname, true); err != nil {
			return nil, err
		}
		log.Println("[INFO] Olric is ready to accept connections")
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
