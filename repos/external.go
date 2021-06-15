package repos

import (
	"net"
	"sync"
	"fmt"

	"github.com/LilithGames/spiracle/config"
)

type External struct {
	Server string
	Addrs  []ExternalAddr
}

func (it External) HostPorts() []string {
	r := make([]string, 0, len(it.Addrs))
	for i := range it.Addrs {
		r = append(r, it.Addrs[i].HostPort)
	}
	return r
}

type ExternalAddr struct {
	HostPort string
	Host     string
	Port     int
	Addr     *net.UDPAddr
}

type ExternelRepo interface {
	Get(server string) (*External, error)
}

type configExternalRepo struct {
	dict *sync.Map
}

func NewConfigExternalRepo(c *config.Config) (ExternelRepo, error){
	dict := &sync.Map{}
	for _, s := range c.RoomProxy.Servers {
		e := &External{Server: s.Name, Addrs: make([]ExternalAddr, 0, len(s.Externals))}
		if len(s.Externals) == 0 {
			return nil, fmt.Errorf("externals is required <%s>", s.Name)
		}
		for _, hostpath := range s.Externals {
			addr, err := net.ResolveUDPAddr("udp4", hostpath)
			if err != nil {
				return nil, fmt.Errorf("invalid externals item <%s> <%s>: %w", s.Name, hostpath, err)
			}
			eaddr := ExternalAddr{
				HostPort: hostpath, 
				Host: addr.IP.String(),
				Port: addr.Port,
				Addr: addr,
			}
			e.Addrs = append(e.Addrs, eaddr)
		}
		dict.Store(s.Name, e)
	}
	return &configExternalRepo{dict}, nil
}

func (it *configExternalRepo) Get(server string) (*External, error) {
	v, ok := it.dict.Load(server)
	if !ok {
		return nil, ErrNotExists
	}
	return v.(*External), nil
}
