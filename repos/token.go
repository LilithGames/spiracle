package repos

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type TokenRepo interface {
	Create(ctx context.Context, opts ...TokenCreationOption) (*Token, error)
	Delete(ctx context.Context, id TToken) error
	Get(ctx context.Context, id TToken) (*Token, error)
}

type Token struct {
	TToken
	Timestamp time.Time
	Expire    time.Time
}

func (it Token) Duration() time.Duration {
	return it.Expire.Sub(it.Timestamp)
}

type tsToken struct {
	dict  sync.Map
	k     int
	hmax  uint32
	lmax  uint32
	lowch chan uint32
}

func NewTsTokenRepo() TokenRepo {
	k := 12
	hmax := uint32(1<<(31-k) - 1)
	lmax := uint32(1<<k - 1)
	max := hmax<<k + lmax
	if max != math.MaxInt32 {
		panic("NewTsTokenRepo max")
	}
	repo := &tsToken{
		k:     k,
		hmax:  hmax,
		lmax:  lmax,
		lowch: make(chan uint32),
	}
	go repo.genLow()
	return repo
}

func (it *tsToken) Create(ctx context.Context, opts ...TokenCreationOption) (*Token, error) {
	o := getTokenCreationOptions(opts...)
	if o.token != 0 {
		ts := time.Now().UTC()
		expire := ts.Add(time.Duration(it.hmax) * time.Second)
		token := &Token{TToken: o.token, Timestamp: ts, Expire: expire}
		if _, loaded := it.dict.LoadOrStore(o.token, token); loaded {
			return nil, fmt.Errorf("create token err: %w", ErrAlreadyExists)
		}
		return token, nil
	}
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("tsToken Create err: %w", ctx.Err())
		default:
			id, ts := it.gen()
			expire := ts.Add(time.Duration(it.hmax) * time.Second)
			token := &Token{TToken: id, Timestamp: ts, Expire: expire}
			if _, loaded := it.dict.LoadOrStore(id, token); loaded {
				continue
			}
			return token, nil
		}
	}
}

func (it *tsToken) Delete(ctx context.Context, id TToken) error {
	it.dict.Delete(id)
	return nil
}

func (it *tsToken) Get(ctx context.Context, id TToken) (*Token, error) {
	v, ok := it.dict.Load(id)
	if !ok {
		return nil, ErrNotExists
	}
	return v.(*Token), nil
}

func (it *tsToken) genLow() {
	var s []int
	for {
		if len(s) == 0 {
			s = rand.Perm(int(it.lmax))
		}
		it.lowch <- uint32(s[len(s)-1] + 1)
		s = s[:len(s)-1]
	}
}

func (it *tsToken) gen() (uint32, time.Time) {
	now := time.Now().UTC()
	h := uint32(now.Unix() % int64(it.hmax))
	l := <-it.lowch
	return h<<it.k + l, now
}

type tokenCreationOptions struct {
	token TToken
}

type TokenCreationOption interface {
	apply(*tokenCreationOptions)
}

type funcTokenCreationOption struct {
	f func(*tokenCreationOptions)
}

func (it *funcTokenCreationOption) apply(o *tokenCreationOptions) {
	it.f(o)
}

func newFuncTokenCreationOption(f func(*tokenCreationOptions)) TokenCreationOption {
	return &funcTokenCreationOption{f: f}
}
func getTokenCreationOptions(opts ...TokenCreationOption) *tokenCreationOptions {
	o := &tokenCreationOptions{}
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

func TokenCreationToken(id TToken) TokenCreationOption {
	return newFuncTokenCreationOption(func(o *tokenCreationOptions) {
		o.token = id
	})
}
