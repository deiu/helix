package helix

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	rdf "github.com/deiu/rdf2go"
	"github.com/stretchr/testify/assert"
)

func Test_GET_RDF(t *testing.T) {
	mime := "text/turtle"
	URI := testServer.URL + "/foo.ttl"
	graph := rdf.NewGraph(URI)
	graph.AddTriple(rdf.NewResource(URI), rdf.NewResource("http://test.com/foo"), rdf.NewLiteral("obj"))
	buf := new(bytes.Buffer)
	graph.Serialize(buf, mime)

	req, err := http.NewRequest("POST", URI, strings.NewReader(buf.String()))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", mime)
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, res.StatusCode)

	req, err = http.NewRequest("GET", URI, nil)
	assert.NoError(t, err)
	req.Header.Add("Accept", mime)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, mime, res.Header.Get("Content-Type"))

	graph = rdf.NewGraph(URI)
	graph.Parse(res.Body, res.Header.Get("Content-Type"))
	res.Body.Close()
	assert.Equal(t, 1, graph.Len())
	assert.NotNil(t, graph.One(rdf.NewResource(URI), rdf.NewResource("http://test.com/foo"), rdf.NewLiteral("obj")))
}

func Test_GET_Static(t *testing.T) {
	req, err := http.NewRequest("GET", testServer.URL+testConfig.StaticPath+"index.html", nil)
	assert.NoError(t, err)
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, "Hello static!", string(body))

	req, err = http.NewRequest("GET", testServer.URL+testConfig.StaticPath+"foo.html", nil)
	assert.NoError(t, err)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
}

func Test_GET_NonRDF(t *testing.T) {
	req, err := http.NewRequest("GET", testServer.URL, nil)
	assert.NoError(t, err)
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	req, err = http.NewRequest("GET", testServer.URL+"/foo", nil)
	assert.NoError(t, err)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
}

func Test_GET_RDFNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", testServer.URL+"/bar", nil)
	assert.NoError(t, err)
	req.Header.Add("Accept", "text/turtle")
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
}

func Test_GET_NotAcceptable(t *testing.T) {
	req, err := http.NewRequest("GET", testServer.URL+"/foo", nil)
	assert.NoError(t, err)
	req.Header.Add("Accept", "text/foo")
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, 406, res.StatusCode)
	assert.NotEmpty(t, body)
}

func Test_HEAD(t *testing.T) {
	req, err := http.NewRequest("HEAD", testServer.URL+"/baz", nil)
	assert.NoError(t, err)
	req.Header.Add("Accept", "text/turtle")
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)

	req, err = http.NewRequest("HEAD", testServer.URL+"/foo.ttl", nil)
	assert.NoError(t, err)
	req.Header.Add("Accept", "text/turtle")
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, 0, len(body))
}

func BenchmarkGET(b *testing.B) {
	mime := "text/turtle"
	URI := testServer.URL + "/benchttl"
	graph := rdf.NewGraph(URI)
	graph.AddTriple(rdf.NewResource(URI), rdf.NewResource("pred"), rdf.NewLiteral("obj"))

	buf := new(bytes.Buffer)
	graph.Serialize(buf, mime)

	req, err := http.NewRequest("POST", URI, strings.NewReader(buf.String()))
	if err != nil {
		b.Fail()
	}
	req.Header.Add("Content-Type", mime)
	_, err = testClient.Do(req)
	if err != nil {
		b.Fail()
	}

	// Run the bench
	e := 0
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", URI, nil)
		if res, _ := testClient.Do(req); res.StatusCode != 200 {
			e++
		}
	}

	// delete resource
	req, err = http.NewRequest("DELETE", URI, nil)
	if err != nil {
		b.Fail()
	}
	_, err = testClient.Do(req)
	if err != nil {
		b.Fail()
	}

	if e > 0 {
		b.Log(fmt.Sprintf("%d/%d failed", e, b.N))
		b.Fail()
	}
}
