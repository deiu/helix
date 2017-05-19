package helix

import (
	"github.com/gocraft/web"
	"log"
)

func (c *Context) GetHandler(w web.ResponseWriter, req *web.Request) {
	var err error
	ctype := ""
	acceptList, _ := conneg(req.Request)
	log.Printf("%+v %+v\n", acceptList, req.Header.Get("Accept"))
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
	graph := c.getGraph(req.RequestURI)
	if graph == nil {
		w.WriteHeader(404)
		return
	}
	graph.Serialize(w, mime)
}
