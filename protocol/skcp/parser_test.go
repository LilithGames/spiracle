package skcp

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)


func TestSKcp(t *testing.T) {
	p := Parser()
	buffer := make([]byte, 26)
	buffer[2] = 0x01
	token, err := p(buffer[1:])
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), token)
}
