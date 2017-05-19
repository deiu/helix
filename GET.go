package helix

import (
	"fmt"
	"github.com/gocraft/web"
)

func (c *Context) GetHandler(w web.ResponseWriter, req *web.Request) {
	logger.Info().Msg("get handler")
	c.Body = "Hello world!"
	fmt.Fprint(w, c.Body)
}
