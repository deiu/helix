package helix

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	chrome = "application/xml,application/xhtml+xml,text/html;q=0.9,text/plain;q=0.8,image/png,*/*;q=0.5"
	rdflib = "application/rdf+xml;q=0.9, application/xhtml+xml;q=0.3, text/xml;q=0.2, application/xml;q=0.2, text/html;q=0.3, text/plain;q=0.1, text/n3;q=1.0, application/x-turtle;q=1, text/turtle;q=1"
)

func mockAccept(accept string) (al AcceptList, err error) {
	req := &http.Request{}
	req.Header = make(http.Header)
	req.Header["Accept"] = []string{accept}
	al, err = conneg(req)
	return
}

func TestNegotiatePicturesOfWebPages(t *testing.T) {
	al, err := mockAccept(chrome)
	if err != nil {
		t.Fatal(err)
	}

	contentType, err := al.Negotiate("text/html", "image/png")
	assert.NoError(t, err)
	assert.Equal(t, "image/png", contentType)
}

func TestNegotiateRDF(t *testing.T) {
	al, err := mockAccept(rdflib)
	if err != nil {
		t.Fatal(err)
	}

	contentType, err := al.Negotiate(serializerMimes...)
	assert.NoError(t, err)
	assert.Equal(t, "text/turtle", contentType)
}

func TestNegotiateFirstMatch(t *testing.T) {
	al, err := mockAccept(chrome)
	assert.NoError(t, err)

	contentType, err := al.Negotiate("text/html", "text/plain", "text/n3")
	assert.NoError(t, err)
	assert.Equal(t, "text/html", contentType)
}

func TestNegotiateSecondMatch(t *testing.T) {
	al, err := mockAccept(chrome)
	assert.NoError(t, err)

	contentType, err := al.Negotiate("text/n3", "text/plain")
	assert.NoError(t, err)
	assert.Equal(t, "text/plain", contentType)
}

func TestNegotiateWildcardMatch(t *testing.T) {
	al, err := mockAccept(chrome)
	assert.NoError(t, err)

	contentType, err := al.Negotiate("text/n3", "application/rdf+xml")
	assert.NoError(t, err)
	assert.Equal(t, "text/n3", contentType)
}

func TestNegotiateInvalidMediaRange(t *testing.T) {
	_, err := mockAccept("something/valid, rubbish, other/valid")
	assert.Error(t, err)
}

func TestNegotiateInvalidParam(t *testing.T) {
	_, err := mockAccept("text/plain; foo")
	assert.Error(t, err)
}

func TestNegotiateEmptyAccept(t *testing.T) {
	al, err := mockAccept("")
	assert.NoError(t, err)

	_, err = al.Negotiate("text/plain")
	assert.Error(t, err)
}

func TestNegotiateNoAlternative(t *testing.T) {
	al, err := mockAccept(chrome)
	assert.NoError(t, err)

	_, err = al.Negotiate()
	assert.Error(t, err)
}

func TestNegotiateStarAccept(t *testing.T) {
	al, err := mockAccept("*")
	assert.NoError(t, err)
	assert.Equal(t, "*/*", al[0].Type+"/"+al[0].SubType)
}
