package helix

import (
	"testing"

	rdf "github.com/deiu/rdf2go"
	"github.com/stretchr/testify/assert"
)

func Test_RDF_CanParse(t *testing.T) {
	ctype := "text/turtle"
	assert.True(t, canParse(ctype))
	ctype = "application/ld+json"
	assert.True(t, canParse(ctype))
	ctype = "application/rdf+xml"
	assert.False(t, canParse(ctype))
}

func Test_RDF_CanSerialize(t *testing.T) {
	ctype := "text/turtle"
	assert.True(t, canSerialize(ctype))
	ctype = "application/ld+json"
	assert.True(t, canSerialize(ctype))
	ctype = "application/rdf+xml"
	assert.False(t, canSerialize(ctype))
}

func Test_RDF_AddRemoveGraph(t *testing.T) {
	c := NewContext()
	URI := "https://example.org"
	graph := rdf.NewGraph(URI)
	c.addGraph(URI, graph)
	assert.Equal(t, URI, c.getGraph(URI).URI())
	c.delGraph(URI)
	assert.Nil(t, c.getGraph(URI))
}
