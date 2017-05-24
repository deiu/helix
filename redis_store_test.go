package helix

import (
	"testing"

	"github.com/albrow/zoom"
	"github.com/stretchr/testify/assert"
)

func Test_Redis_InitEmptyURL(t *testing.T) {
	var testPool *zoom.Pool
	_, err := initRedisPool("", testPool)
	assert.Error(t, err)
}

func Test_Redis_URL(t *testing.T) {
	var testPool *zoom.Pool
	_, err := initRedisPool("localhost:1234", testPool)
	assert.NoError(t, err)
}
