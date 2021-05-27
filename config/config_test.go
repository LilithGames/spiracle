package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c, err := Load("../config.yaml")
	assert.Nil(t, err)
	fmt.Printf("%+v\n", c)
	assert.Equal(t, 2, len(c.RoomProxy.Servers))
}
