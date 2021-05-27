package multiplex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	p := Parser()
	token, err := p([]byte{0x01})
	assert.Nil(t, err)
	assert.Equal(t, byte(0x01), token.(byte))
	token, err = p([]byte{'x'})
	assert.Nil(t, err)
	assert.Equal(t, byte('x'), token.(byte))
	token, err = p([]byte{0x00})
	assert.NotNil(t, err)
}
