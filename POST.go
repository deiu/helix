package helix

import (
	rdf "github.com/deiu/rdf2go"
	"github.com/gocraft/web"
)

func (c *Context) PostHandler(w web.ResponseWriter, req *web.Request) {
	if canParse(req.Header.Get("Content-Type")) {
		c.postRDF(w, req)
		return
	}
	w.WriteHeader(400)
}

func (c *Context) postRDF(w web.ResponseWriter, req *web.Request) {
	graph := rdf.NewGraph(req.RequestURI)
	graph.Parse(req.Body, req.Header.Get("Content-Type"))
	if graph.Len() == 0 {
		w.WriteHeader(400)
		w.Write([]byte("Empty request body"))
		return
	}
	c.addGraph(req.RequestURI, graph)
}
