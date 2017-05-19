package helix

import (
	"github.com/gocraft/web"
	"github.com/rs/zerolog"
	"os"
)

func (c *Context) RequestLogger(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	logger = zerolog.New(os.Stderr).With().
		Timestamp().
		Str("Method", r.Method).
		Str("Path", r.Request.URL.String()).
		Logger()

	next(w, r)
}
