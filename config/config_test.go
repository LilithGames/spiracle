package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Setenv("CONFIGOR_ROOMPROXY_ENABLE", "false")
	c, err := Load("../config.yaml")
	assert.Nil(t, err)
	fmt.Printf("%+v\n", c)
	assert.Equal(t, 2, len(c.RoomProxy.Servers))
	println(c.RoomProxy.Enable)
}
