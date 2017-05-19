package helix

import (
	"github.com/gocraft/web"
	"strings"
)

func (c *Context) OptionsHandler(w web.ResponseWriter, r *web.Request, methods []string) {
	origin := r.Request.Header.Get("Origin")
	if len(origin) > 0 {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	if len(origin) < 1 {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	crh := r.Request.Header.Get("Access-Control-Request-Headers") // CORS preflight only
	if len(crh) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", crh)
	}
	crm := r.Request.Header.Get("Access-Control-Request-Method") // CORS preflight only
	if len(crm) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", crm)
	} else {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methodsAll, ", "))
	}

	// w.Header().Add("Access-Control-Allow-Methods", strings.Join(methods, ", "))
}
