package kcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKcpBasic(t *testing.T) {
	p := Parser()
	buffer := make([]byte, 25)
	buffer[1] = 0x01
	token, err := p(buffer[1:])
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), token.(uint32))
}
