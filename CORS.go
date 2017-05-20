package helix

import (
	"github.com/gocraft/web"
	"strings"
)

func (c *Context) CORSMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
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

	w.Header().Set("Accept-Post", strings.Join(rdfMimes, ", "))

	if c.Config.HSTS {
		w.Header().Set("Strict-Transport-Security", "max-age=63072000")
	}

	next(w, r)
}
