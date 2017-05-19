package helix

import (
	"fmt"
	"github.com/gocraft/web"
)

func (c *Context) GetHandler(w web.ResponseWriter, req *web.Request) {
	// c.g.Graph = "https://example.org/foo"
	// g.Subject = "https://example.org/foo#this"
	// g.Predicate = "a"
	// g.Object = "http://xmlns.com/foaf/0.1/PersonalProfileDocument"
	// g.Current = true
	logger.Info().Msg("get handler")
	c.Body = "Hello world!"
	fmt.Fprint(w, c.Body)
}
