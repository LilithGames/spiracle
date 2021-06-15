package repos

import (
	"testing"
	"fmt"

	"github.com/LilithGames/spiracle/config"
	"github.com/stretchr/testify/assert"
)

func TestExternalRepo(t *testing.T) {
	c, err := config.Load("../config.yaml")
	assert.Nil(t, err)
	erepo, err := NewConfigExternalRepo(c)
	assert.Nil(t, err)
	e, err := erepo.Get("dev")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(e.Addrs))
	fmt.Printf("%+v\n", e)
}
