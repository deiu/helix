package helix

import (
	"github.com/gocraft/web"
)

func (c *Context) GetHandler(w web.ResponseWriter, req *web.Request) {
	var err error
	ctype := ""
	acceptList, _ := conneg(req.Request)
	if len(acceptList) > 0 && acceptList[0].SubType != "*" {
		ctype, err = acceptList.Negotiate(serializerMimes...)
		if err != nil {
			w.WriteHeader(406)
			w.Write([]byte("HTTP 406 - Accept type not acceptable: " + err.Error()))
			return
		}
		logger.Info().Str("Accept", ctype).Msg("")
	}

	if canSerialize(ctype) {
		c.getRDF(w, req, ctype)
		return
	}
	// w.WriteHeader(400)
}

func (c *Context) getRDF(w web.ResponseWriter, req *web.Request, mime string) {
	graph, err := c.getGraph(req.RequestURI)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
	graph.Serialize(w, mime)
}
