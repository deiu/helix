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

func Test_Negotiate_PicturesOfWebPages(t *testing.T) {
	al, err := mockAccept(chrome)
	if err != nil {
		t.Fatal(err)
	}

	contentType, err := al.Negotiate("text/html", "image/png")
	assert.NoError(t, err)
	assert.Equal(t, "image/png", contentType)
}

func Test_Negotiate_RDF(t *testing.T) {
	al, err := mockAccept(rdflib)
	if err != nil {
		t.Fatal(err)
	}

	contentType, err := al.Negotiate(rdfMimes...)
	assert.NoError(t, err)
	assert.Equal(t, "text/turtle", contentType)
}

func Test_Negotiate_FirstMatch(t *testing.T) {
	al, err := mockAccept(chrome)
	assert.NoError(t, err)

	contentType, err := al.Negotiate("text/html", "text/plain", "text/n3")
	assert.NoError(t, err)
	assert.Equal(t, "text/html", contentType)
}

func Test_Negotiate_SecondMatch(t *testing.T) {
	al, err := mockAccept(chrome)
	assert.NoError(t, err)

	contentType, err := al.Negotiate("text/n3", "text/plain")
	assert.NoError(t, err)
	assert.Equal(t, "text/plain", contentType)
}

func Test_Negotiate_WildcardMatch(t *testing.T) {
	al, err := mockAccept(chrome)
	assert.NoError(t, err)

	contentType, err := al.Negotiate("text/n3", "application/rdf+xml")
	assert.NoError(t, err)
	assert.Equal(t, "text/n3", contentType)
}

func Test_Negotiate_SubType(t *testing.T) {
	al, err := mockAccept("text/turtle, application/*")
	assert.NoError(t, err)

	contentType, err := al.Negotiate("foo/bar", "application/ld+json")
	assert.NoError(t, err)
	assert.Equal(t, "application/ld+json", contentType)
}

func Test_Negotiate_InvalidMediaRange(t *testing.T) {
	_, err := mockAccept("something/valid, fail, other/valid")
	assert.Error(t, err)
}

func Test_Negotiate_Invalid_Param(t *testing.T) {
	_, err := mockAccept("text/plain; foo")
	assert.Error(t, err)
}

func Test_Negotiate_OtherParam(t *testing.T) {
	_, err := mockAccept("text/plain;foo=bar")
	assert.NoError(t, err)
}

func Test_Negotiate_EmptyAccept(t *testing.T) {
	al, err := mockAccept("")
	assert.NoError(t, err)

	_, err = al.Negotiate("text/plain")
	assert.Error(t, err)
}

func Test_Negotiate_NoAlternative(t *testing.T) {
	al, err := mockAccept(chrome)
	assert.NoError(t, err)

	_, err = al.Negotiate()
	assert.Error(t, err)
}

func Test_Negotiate_StarAccept(t *testing.T) {
	al, err := mockAccept("*")
	assert.NoError(t, err)
	assert.Equal(t, "*/*", al[0].Type+"/"+al[0].SubType)

	al, err = mockAccept("text/*")
	assert.NoError(t, err)
	assert.Equal(t, "text/*", al[0].Type+"/"+al[0].SubType)
}

func Test_Negotiate_Sorter(t *testing.T) {
	accept := []Accept{}
	a := Accept{
		Type:    "text",
		SubType: "*",
		Q:       float32(3),
	}
	accept = append(accept, a)

	a = Accept{
		Type:    "*",
		SubType: "text",
		Q:       float32(5),
	}
	accept = append(accept, a)

	a = Accept{
		Type:    "*",
		SubType: "text",
	}
	accept = append(accept, a)

	sorter := acceptSorter(accept)
	assert.True(t, sorter.Less(0, 1))
	assert.True(t, sorter.Less(2, 0))
}
