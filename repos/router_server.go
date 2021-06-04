package repos

import (
	"fmt"
	"regexp"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/query"
)

type routerRepo struct {
	db   *olric.Olric
	dmap *olric.DMap
}

func NewRouterRepo(db *olric.Olric) (RouterRepo, error) {
	dmap, err := db.NewDMap("storage.routers")
	if err != nil {
		return nil, err
	}
	return &routerRepo{
		dmap: dmap,
		db:   db,
	}, nil
}

func (it *routerRepo) name(scope string, name string) string {
	return scope + "." + name
}

func (it *routerRepo) Create(record *RouterRecord, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	key := it.name(o.scope, str(record.Token))
	return it.dmap.PutIf(key, record, olric.IfNotFound)
}

func (it *routerRepo) Update(record *RouterRecord, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	key := it.name(o.scope, str(record.Token))
	return it.dmap.PutIf(key, record, olric.IfFound)
}

func (it *routerRepo) CreateOrUpdate(record *RouterRecord, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	key := it.name(o.scope, str(record.Token))
	return it.dmap.Put(key, record)
}

func (it *routerRepo) Delete(token TToken, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	key := it.name(o.scope, str(token))
	return it.dmap.Delete(key)
}

func (it *routerRepo) Get(token TToken, opts ...RouterOption) (*RouterRecord, error) {
	o := getRouterOptions(opts...)
	key := it.name(o.scope, str(token))
	r, err := it.dmap.Get(key)
	if err != nil {
		return nil, err
	}
	return r.(*RouterRecord), nil
}

func (it *routerRepo) List(f func(*RouterRecord) bool, opts ...RouterOption) error {
	o := getRouterOptions(opts...)
	_ = o
	cur, err := it.dmap.Query(query.M{"$onKey": query.M{"$regexMatch": fmt.Sprintf("^%s\\.", regexp.QuoteMeta(o.scope))}})
	if err != nil {
		return err
	}
	defer cur.Close()
	return cur.Range(func(key string, value interface{}) bool {
		return f(value.(*RouterRecord))
	})
}
