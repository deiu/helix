package helix

import (
	rdf "github.com/deiu/rdf2go"
	"github.com/gocraft/web"
)

func (c *Context) PostHandler(w web.ResponseWriter, req *web.Request) {
	ctype := req.Header.Get("Content-Type")
	logger.Info().Str("Content-Type", ctype).Msg("")
	if canParse(ctype) {
		c.postRDF(w, req)
		return
	}
	w.WriteHeader(400)
}

func (c *Context) postRDF(w web.ResponseWriter, req *web.Request) {
	URI := absoluteURI(req.Request)
	graph := rdf.NewGraph(URI, true)
	graph.Parse(req.Body, req.Header.Get("Content-Type"))
	if graph.Len() == 0 {
		w.WriteHeader(400)
		w.Write([]byte("Empty request body"))
		return
	}
	_, err := c.getGraph(URI)
	if err == nil {
		w.WriteHeader(409)
		w.Write([]byte("Cannot create new graph if it aready exists"))
		return
	}
	c.addGraph(URI, graph)

	// add ETag
	w.Header().Add("ETag", newETag([]byte(graph.String())))

	w.WriteHeader(201)
}
