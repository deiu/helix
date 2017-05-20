package helix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RDF_NewETag(t *testing.T) {
	data := []byte("test")
	assert.Equal(t, "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", newETag(data))
}

func Test_ETag_Match(t *testing.T) {
	goodHeader := "12345"
	badHeader := "123"
	empty := ""
	star := "*"
	etag := "12345"

	assert.False(t, ETagMatch(badHeader, etag))
	assert.True(t, ETagMatch("", ""))
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
	assert.True(t, ETagNoneMatch("", ""))
	assert.True(t, ETagNoneMatch(empty, etag))
	assert.False(t, ETagNoneMatch(goodHeader, etag))
	assert.False(t, ETagNoneMatch(star, etag))
}
