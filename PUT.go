package helix

import (
	"fmt"
	"github.com/gocraft/web"
)

func (c *Context) PutHandler(w web.ResponseWriter, req *web.Request) {
	fmt.Fprint(w, c.Body)
}
