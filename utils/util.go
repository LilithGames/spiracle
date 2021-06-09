package utils

import (
	"net"
	"net/http"
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

func Readyz(done chan struct{}) healthz.Checker {
	once := &sync.Once{}
	return func(_ *http.Request) error {
		once.Do(func() {
			close(done)
		})
		return nil
	}
}

func ResolveUDP4Addr(addr string) *net.UDPAddr {
	result, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil
	}
	return result
}
