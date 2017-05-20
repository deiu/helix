package helix

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	rdf "github.com/deiu/rdf2go"
	"github.com/stretchr/testify/assert"
)

func Test_DELETE_NonExistent(t *testing.T) {
	req, err := http.NewRequest("DELETE", testServer.URL+"/foo", nil)
	assert.NoError(t, err)
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
}

func Test_DELETE_RDF(t *testing.T) {
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

	req, err = http.NewRequest("DELETE", URI, nil)
	assert.NoError(t, err)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	req, err = http.NewRequest("GET", URI, nil)
	assert.NoError(t, err)
	req.Header.Add("Accept", mime)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
}
