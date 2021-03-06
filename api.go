// Package api is a helper that simplifies the process of REST APIs bindings creation in Go.
// Rather than composing URLs and HTTP requests by hand, one can use the api.Request method in order to
// automatically create such a request. The use case may be as following:
//   svc, _ := api.New("http://example.com")
//   args := url.Values{}
//   args.Set("filter", "1")
//   args.Set("price", "200")
//   req, _ := svc.Request(api.GET, "/categories/1", args)
//
//   // URL is now http://example.com/categories/1?filter=1&price=200
//
//   var cli http.Client
//   resp, err := cli.Do(req)
//
// In the case of POST, the arguments will be presented in the Body of request:
//
//   req, _ := svc.Request(api.POST, "/categories/1", args)
//
//   // URL is now http://example.com/categories/1
//   // Body is now filter=1&price=200
//   // Header is now has Content-Type: application/x-www-form-urlencoded
//
//   var cli http.Client
//   resp, err := cli.Do(req)
package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

// Method represents an HTTP method.
type Method int

const (
	GET Method = iota
	POST
	HEAD
	PUT
	DELETE
	PATCH
)

func (m Method) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case HEAD:
		return "HEAD"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case PATCH:
		return "PATCH"
	default:
		return "GET"
	}
}

// Api represents a REST API connection.
type Api struct {
	// BaseURI is the base URI of an API.
	BaseURI *url.URL
	// Header is a custom header that will be used for communtication with API (e.g. Authorization).
	Header http.Header
}

// New creates a new api instance with given base uri.
func New(uri string) (a *Api, err error) {
	a = &Api{}
	a.BaseURI, err = url.ParseRequestURI(uri)
	return
}

// MustNew is like New, but panics if any error has occured.
func MustNew(uri string) *Api {
	a, err := New(uri)
	if err != nil {
		panic(err)
	}
	return a
}

// Request creates an http request instance properly initialized with the given parameters.
// In a special case for the POST method it will create a body buffer,
// in other cases it will just store the parameters in the URL.
func (a *Api) Request(method Method, resource string, args url.Values) (req *http.Request, err error) {
	u := *a.BaseURI
	u.Path = path.Join(u.Path, resource)

	switch method {
	case GET, HEAD, PUT, DELETE, PATCH:
		u.RawQuery = args.Encode()
		if req, err = http.NewRequest(method.String(), u.String(), nil); err != nil {
			return
		}
		for k := range a.Header {
			req.Header.Add(k, a.Header.Get(k))
		}
	case POST:
		data := args.Encode()
		if req, err = http.NewRequest(method.String(), u.String(), bytes.NewBufferString(data)); err != nil {
			return
		}
		for k := range a.Header {
			req.Header.Add(k, a.Header.Get(k))
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Content-Length", strconv.Itoa(len(data)))
	default:
		return nil, fmt.Errorf("api: unknown method: %d", method)
	}

	return req, nil
}

func (a *Api) RequestBytes(method Method, resource string, contentType string, data []byte) (req *http.Request, err error) {
	u := *a.BaseURI
	u.Path = path.Join(u.Path, resource)
	if req, err = http.NewRequest(method.String(), u.String(), bytes.NewReader(data)); err != nil {
		return
	}
	for k := range a.Header {
		req.Header.Add(k, a.Header.Get(k))
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))
	return
}
