package helix

import (
	"fmt"
	"github.com/gocraft/web"
	"net/http"
)

func (c *Context) RootHandler(w web.ResponseWriter, req *web.Request) {
	logger.Info().Msg("In root")
	c.Body = "Hello world from root"
	fmt.Fprint(w, c.Body)
}

func (c *Context) NotFound(w web.ResponseWriter, r *web.Request) {
	w.WriteHeader(http.StatusNotFound) // You probably want to return 404. But you can also redirect or do whatever you want.
	fmt.Fprintf(w, "404 - Not Found")  // Render you own HTML or something!
}
