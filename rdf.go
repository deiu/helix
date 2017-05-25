package helix

import (
	"errors"
	rdf "github.com/deiu/rdf2go"
)

var mimeParser = map[string]string{
	"text/turtle":         "turtle",
	"application/ld+json": "jsonld",
	// "application/sparql-update": "internal",
}

var mimeSerializer = map[string]string{
	"text/turtle":         "turtle",
	"application/ld+json": "jsonld",
}

var rdfMimes = []string{
	"text/turtle",
	"application/ld+json",
}

type Graph struct {
	*rdf.Graph
	Etag string
}

func NewGraph() *Graph {
	return &Graph{}
}

func canParse(ctype string) bool {
	if len(mimeParser[ctype]) > 0 {
		return true
	}
	return false
}

func canSerialize(ctype string) bool {
	if len(mimeSerializer[ctype]) > 0 {
		return true
	}
	return false
}

func (c *Context) addGraph(URI string, graph *Graph) {
	c.Store[URI] = graph
}

func (c *Context) getGraph(URI string) (*Graph, error) {
	if c.Store[URI] == nil {
		return nil, errors.New("Cannot find graph that matches URI: " + URI)
	}
	return c.Store[URI], nil
}

func (c *Context) delGraph(URI string) error {
	if c.Store[URI] == nil {
		return errors.New("Cannot delete graph that matches URI: " + URI)
	}
	delete(c.Store, URI)
	return nil
}
