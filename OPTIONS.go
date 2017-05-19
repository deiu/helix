package helix

import (
	"github.com/gocraft/web"
	"strings"
)

func (c *Context) OptionsHandler(w web.ResponseWriter, r *web.Request, methods []string) {
	w.Header().Add("Access-Control-Allow-Methods", strings.Join(methods, ", "))
	w.Header().Add("Access-Control-Allow-Origin", "*")
}
