package helix

import (
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

var serializerMimes = []string{
	"text/turtle",
	"application/ld+json",
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

func (c *Context) addGraph(URI string, graph *rdf.Graph) {
	c.Store[URI] = graph
}

func (c *Context) getGraph(URI string) *rdf.Graph {
	return c.Store[URI]
}

func (c *Context) delGraph(URI string) {
	delete(c.Store, URI)
}
