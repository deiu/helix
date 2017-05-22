package helix

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	rdf "github.com/deiu/rdf2go"
	"github.com/stretchr/testify/assert"
)

func Test_POST_OtherMime(t *testing.T) {
	URI := testServer.URL + "/foo"
	req, err := http.NewRequest("POST", URI, strings.NewReader("foo"))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "text/plain")
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, strings.Join(rdfMimes, ", "), res.Header.Get("Accept-Post"))
}

func Test_POST_TurtleEmpty(t *testing.T) {
	mime := "text/turtle"
	URI := testServer.URL + "/foo"

	req, err := http.NewRequest("POST", URI, nil)
	assert.NoError(t, err)
	req.Header.Add("Content-Type", mime)
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

func Test_POST_Turtle(t *testing.T) {
	mime := "text/turtle"
	URI := testServer.URL + "/foo"
	graph := rdf.NewGraph(URI)
	graph.AddTriple(rdf.NewResource(URI), rdf.NewResource("pred"), rdf.NewLiteral("obj"))

	buf := new(bytes.Buffer)
	graph.Serialize(buf, mime)

	req, err := http.NewRequest("POST", URI, strings.NewReader(buf.String()))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", mime)
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, res.StatusCode)
	etag := res.Header.Get("Etag")
	assert.NotEmpty(t, etag)

	req, err = http.NewRequest("POST", URI, strings.NewReader(buf.String()))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", mime)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 409, res.StatusCode)

	req, err = http.NewRequest("GET", URI, nil)
	assert.NoError(t, err)
	req.Header.Add("Accept", mime)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, etag, res.Header.Get("Etag"))
}
