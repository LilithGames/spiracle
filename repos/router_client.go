package repos

import (
	"fmt"
	"regexp"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/client"
	"github.com/buraksezer/olric/query"
)

type routerClientRepo struct {
	db   *client.Client
	dmap *client.DMap
}

func NewClientRouterRepo(db *client.Client) RouterRepo {
	dmap := db.NewDMap("storage.routers")
	return &routerClientRepo{
		dmap: dmap,
		db:   db,
	}
}

func (it *routerClientRepo) name(scope string, name string) string {
	return scope + ":" + name
}

func (it *routerClientRepo) Create(record *RouterRecord, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	key := it.name(o.scope, str(record.Token))
	return it.dmap.PutIf(key, record, olric.IfNotFound)
}

func (it *routerClientRepo) Update(record *RouterRecord, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	key := it.name(o.scope, str(record.Token))
	return it.dmap.PutIf(key, record, olric.IfFound)
}

func (it *routerClientRepo) CreateOrUpdate(record *RouterRecord, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	key := it.name(o.scope, str(record.Token))
	return it.dmap.Put(key, record)
}

func (it *routerClientRepo) Delete(token TToken, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	key := it.name(o.scope, str(token))
	return it.dmap.Delete(key)
}

func (it *routerClientRepo) Get(token TToken, opts ...RouterOption) (*RouterRecord, error) {
	o := getRouterOptions(opts...)
	key := it.name(o.scope, str(token))
	r, err := it.dmap.Get(key)
	if err != nil {
		return nil, err
	}
	return r.(*RouterRecord), nil
}

func (it *routerClientRepo) List(f func(*RouterRecord) bool, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	cur, err := it.dmap.Query(query.M{"$onKey": query.M{"$regexMatch": fmt.Sprintf("^%s\\.", regexp.QuoteMeta(o.scope))}})
	if err != nil {
		return err
	}
	defer cur.Close()
	return cur.Range(func(key string, value interface{}) bool {
		return f(value.(*RouterRecord))
	})
}

