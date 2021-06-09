package repos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	repo := NewTsTokenRepo()
	token, err := repo.Create(context.TODO())
	assert.Nil(t, err)
	fmt.Printf("%+v\n", token)
	_, err = repo.Create(context.TODO(), TokenCreationToken(token.TToken))
	assert.Equal(t, err, ErrAlreadyExists)
}

func TestTokenQPS(t *testing.T) {
	repo := NewTsTokenRepo()
    for i := 0; i < 8000; i++ {
		_, err := repo.Create(context.TODO())
		assert.Nil(t, err)
    }
}
