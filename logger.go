package helix

import (
	"os"

	"github.com/gocraft/web"
	"github.com/rs/zerolog"
)

func (c *Context) RequestLogger(w web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	logger = zerolog.New(os.Stderr).With().
		Timestamp().
		Str("Method", req.Method).
		Str("Path", req.Request.URL.String()).
		Str("User", c.User).
		Logger()

	next(w, req)
}
