API
---
[![GoDoc](https://godoc.org/github.com/xlab/api?status.svg)](https://godoc.org/github.com/xlab/api)

#### Installation:
```
go get github.com/xlab/api
```

#### Use case:
```go
svc, _ := api.New("http://example.com")
args := url.Values{}
args.Set("filter", "1")
args.Set("price", "200")
req, _ := svc.Request(api.POST, "/categories/1", args)

// req.URL is now http://example.com/categories/1
// req.Body is now filter=1&price=200
// req.Header now has Content-Type: application/x-www-form-urlencoded

var cli http.Client
resp, err := cli.Do(req)
```

#### License

MIT
