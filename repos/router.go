package repos

import (
	"net"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/query"
)

type TToken = uint32

type RouterRecord struct {
	Token  TToken
	Addr   *net.UDPAddr
	RoomId string
}

type RouterRepo interface {
	Create(record *RouterRecord) error
	Update(record *RouterRecord) error
	CreateOrUpdate(record *RouterRecord) error
	Delete(token TToken) error
	Get(token TToken) (*RouterRecord, error)
	List(f func(*RouterRecord) bool) error
}

type routerRepo struct {
	name string
	db   *olric.Olric
	dmap *olric.DMap
}

func NewRouterRepo(scope string, db *olric.Olric) (RouterRepo, error) {
	if scope == "" {
		scope = "$global"
	}
	name := "storage." + scope + ".router"
	dmap, err := db.NewDMap(name)
	if err != nil {
		return nil, err
	}
	return &routerRepo{
		name: name,
		dmap: dmap,
		db:   db,
	}, nil
}

func (it *routerRepo) Create(record *RouterRecord) error {
	return it.dmap.PutIf(str(record.Token), record, olric.IfNotFound)
}

func (it *routerRepo) Update(record *RouterRecord) error {
	return it.dmap.PutIf(str(record.Token), record, olric.IfFound)
}

func (it *routerRepo) CreateOrUpdate(record *RouterRecord) error {
	return it.dmap.Put(str(record.Token), record)
}

func (it *routerRepo) Delete(token TToken) error {
	return it.dmap.Delete(str(token))
}

func (it *routerRepo) Get(token TToken) (*RouterRecord, error) {
	r, err := it.dmap.Get(str(token))
	if err != nil {
		return nil, err
	}
	return r.(*RouterRecord), nil
}

func (it *routerRepo) List(f func(*RouterRecord) bool) error {
	cur, err := it.dmap.Query(query.M{"$onKey": query.M{"$regexMatch": ""}})
	if err != nil {
		return err
	}
	defer cur.Close()
	return cur.Range(func(key string, value interface{}) bool {
		return f(value.(*RouterRecord))
	})
}
