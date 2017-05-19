package helix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ETag_Match(t *testing.T) {
	goodHeader := "12345"
	badHeader := "123"
	empty := ""
	star := "*"
	etag := "12345"

	assert.False(t, ETagMatch(badHeader, etag))
	assert.True(t, ETagMatch(empty, etag))
	assert.True(t, ETagMatch(goodHeader, etag))
	assert.True(t, ETagMatch(star, etag))
}

func Test_ETag_NoneMatch(t *testing.T) {
	goodHeader := "12345"
	badHeader := "123"
	empty := ""
	star := "*"
	etag := "12345"

	assert.True(t, ETagNoneMatch(badHeader, etag))
	assert.True(t, ETagNoneMatch(empty, etag))
	assert.False(t, ETagNoneMatch(goodHeader, etag))
	assert.False(t, ETagNoneMatch(star, etag))
}
