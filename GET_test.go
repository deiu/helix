package helix

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, 200, res.StatusCode)
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

func BenchmarkGET(b *testing.B) {
	e := 0
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", testServer.URL, nil)
		if res, _ := testClient.Do(req); res.StatusCode != 200 {
			e++
		}
	}
	if e > 0 {
		b.Log(fmt.Sprintf("%d/%d failed", e, b.N))
		b.Fail()
	}
}
