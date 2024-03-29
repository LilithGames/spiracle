package heartbeat

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	p := Parser()
	token, err := p([]byte("1234 xxxxx"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(1234), token)
	token, err = p([]byte("1234xxxxx"))
	assert.NotNil(t, err)
	token, err = p([]byte("hello xxxxx"))
	assert.NotNil(t, err)
}
