package helix

import (
	"fmt"
	"github.com/gocraft/web"
)

func (c *Context) RootHandler(w web.ResponseWriter, req *web.Request) {
	logger.Info().Msg("In root")
	fmt.Fprint(w, "Hello world from root")
}
