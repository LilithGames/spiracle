package controllers

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	d := map[string]struct{}{
		"hello": struct{}{},
		"hello1": struct{}{},
	}
	ks := keys(d)
	assert.Contains(t, ks, "hello", "hello1")
}
