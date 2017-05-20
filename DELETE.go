package helix

import (
	"github.com/gocraft/web"
)

func (c *Context) DeleteHandler(w web.ResponseWriter, req *web.Request) {
	err := c.delGraph(req.RequestURI)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
}
