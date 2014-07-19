package api

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestGet(t *testing.T) {
	a, err := New("http://example.com")
	if !assert.NoError(t, err) {
		return
	}
	args := url.Values{}
	args.Add("filter", "1")
	args.Add("price", "200")
	req, err := a.Request(GET, "/categories/1", args)
	assert.NoError(t, err)
	expURL := "http://example.com/categories/1?filter=1&price=200"
	assert.Equal(t, expURL, req.URL.String())
}

func TestRequestPost(t *testing.T) {
	a, err := New("http://example.com")
	if !assert.NoError(t, err) {
		return
	}
	args := url.Values{}
	args.Add("filter", "1")
	args.Add("price", "200")
	a.Header = http.Header{}
	a.Header.Add("foo", "bar")
	req, err := a.Request(POST, "/categories/1", args)
	assert.NoError(t, err)
	defer req.Body.Close()
	expURL := "http://example.com/categories/1"
	assert.Equal(t, expURL, req.URL.String())
	expBody := "filter=1&price=200"
	buf, _ := ioutil.ReadAll(req.Body)
	assert.Equal(t, expBody, string(buf))
}

func TestRequestHeaders(t *testing.T) {
	a, err := New("http://example.com")
	if !assert.NoError(t, err) {
		return
	}
	args := url.Values{}
	a.Header = http.Header{}
	a.Header.Add("foo", "bar")
	req, err := a.Request(GET, "/categories/1", args)
	assert.NoError(t, err)
	expHeader := http.Header{
		"Foo": []string{"bar"},
	}
	assert.Equal(t, expHeader, req.Header)
}

func TestRequestErrors(t *testing.T) {
	a, err := New("example.com")
	assert.Error(t, err)
	a, err = New("http://example.com")
	_, err = a.Request(Method(10), "", nil)
	assert.Error(t, err)
}
