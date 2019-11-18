# go-restclient [![GoDoc][doc-img]][doc]

A higher-order type that simplifies calling REST APIs. It wraps the standard http and json packages to provide the following features:

- JSON encoding of request entity
- JSON parsing of response entity
- conversion of non-2xx responses into a structured error type including response body extraction
- timeout management
- easier request URL building given a client-wide base URL and query values
- auto-closing of the response body
- interceptors to inject authentication tokens, etc

## Installation

`go get github.com/racker/go-restclient`

## Quick example

```go
package main

import (
	"fmt"
	"github.com/racker/go-restclient"
)

func main() {
	client := restclient.New()
	client.SetBaseUrl("http://your.own.domain")

	type MsgHolder struct {
		Msg string
	}

	req := &MsgHolder{Msg:"hello"}
	var resp MsgHolder

	client.Exchange("POST", "/ping", nil,
		restclient.NewJsonEntity(req), restclient.NewJsonEntity(&resp))
    
  	fmt.Println(resp.Msg)
}
```

[doc-img]: https://godoc.org/github.com/racker/go-restclient?status.svg
[doc]: https://godoc.org/github.com/racker/go-restclient
