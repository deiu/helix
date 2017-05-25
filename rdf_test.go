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
	var err error
	c := NewContext()
	URI := "https://example.org"
	g := rdf.NewGraph(URI)
	graph := NewGraph()
	graph.Graph = g
	c.addGraph(URI, graph)
	graph, err = c.getGraph(URI)
	assert.NoError(t, err)
	assert.Equal(t, URI, graph.Graph.URI())
	err = c.delGraph(URI)
	assert.NoError(t, err)
	graph, err = c.getGraph(URI)
	assert.Error(t, err)
	assert.Nil(t, graph)
}
