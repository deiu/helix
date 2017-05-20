package helix

import (
	"github.com/gocraft/web"
)

func (c *Context) DeleteHandler(w web.ResponseWriter, req *web.Request) {
	URI := absoluteURI(req.Request)
	err := c.delGraph(URI)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
}
