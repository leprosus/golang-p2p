# Golang TCP simple client and server

## Import

```go
import "github.com/leprosus/golang-p2p"
```

## Create new TCP

```go
tcp := p2p.NewTCP("localhost", 8080)
```

## Server example

```go
package main

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
)

func main() {
	var c int

	tcp := p2p.NewTCP("localhost", 8080)
	err := tcp.Handle(func(req *p2p.Request, res *p2p.Response) {
		fmt.Println(">", string(req.Body))

		c++
		err := res.Send([]byte(fmt.Sprintf("buy %d", c)))
		if err != nil {
			panic(err)
		}
	})
	if err != nil {
		panic(err)
	}
}
```

## Client example

```go
package main

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
)

func main() {
	for i := 0; i < 10; i++ {
		tcp := p2p.NewTCP("localhost", 8080)
		res, err := tcp.Send([]byte("hello"))
		if err != nil {
			panic(err)
		}

		fmt.Println("<", string(res))
	}
}
```

## Running

If you run the server and the client separately then you see:

* in the server stdout:

```
> hello
> hello
> hello
> hello
> hello
> hello
> hello
> hello
> hello
> hello
```

* in the client stdout:

```
< buy 1
< buy 2
< buy 3
< buy 4
< buy 5
< buy 6
< buy 7
< buy 8
< buy 9
< buy 10
```

## List all methods

### TCP Initialization

* p2p.NewTCP(host, port) - creates TCP connection
* p2p.SetTimeout(timeout) - sets timeout on a communication
* p2p.GetTimeout() - returns current timeout on a communication (30 seconds is default value)
* p2p.SetRequestLimit(limit) - sets request body limit
* p2p.GetRequestLimit() - returns request body limit (1024 bytes is default value)
* p2p.Close() - closes all connections

### Server

* tcp.Handle(handler) - sets `func(req *Request, res *Response)` to handle requests and send responses
* req.Body - the request field contains sent data of a client
* res.Send(bs) - sends bytes back to client

### Client

* tcp.Send(bs) - sends bytes to server
