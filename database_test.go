package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenGetPut(t *testing.T) {
	assert.Nil(t, Open())
	defer Close()

	assert.Nil(t, Put([]byte("key"), []byte("value")))

	value, err := Get([]byte("key"))
	assert.Equal(t, []byte("value"), value)
	assert.Nil(t, err)
}
