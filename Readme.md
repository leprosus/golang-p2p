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
	"context"
	"fmt"
	"log"

	p2p "github.com/leprosus/golang-p2p"
)

type Hello struct {
	Text string
}

type Buy struct {
	Text string
}

func main() {
	tcp := p2p.NewTCP("localhost", 8080)
	settings := p2p.NewServerSettings()

	server := p2p.NewServer(tcp, settings)

	server.SetHandle("dialog", func(ctx context.Context, req p2p.Request) (res p2p.Response, err error) {
		hello := Hello{}
		err = req.GetGob(&hello)
		if err != nil {
			return
		}

		fmt.Printf("> Hello: %s\n", hello.Text)

		buy := Buy{Text: hello.Text}
		err = res.SetGob(buy)

		return
	})

	err := server.Serve()
	if err != nil {
		log.Panicln(err)
	}
}
```

## Client example

```go
package main

import (
	"fmt"
	"log"

	p2p "github.com/leprosus/golang-p2p"
)

type Hello struct {
	Text string
}

type Buy struct {
	Text string
}

func main() {
	tcp := p2p.NewTCP("localhost", 8080)
	settings := p2p.NewClientSettings()

	client := p2p.NewClient(tcp, settings)

	for i := 0; i < 10; i++ {
		hello := Hello{Text: fmt.Sprintf("User #%d", i+1)}

		req := p2p.Request{}
		err := req.SetGob(hello)
		if err != nil {
			log.Panicln(err)
		}

		var res p2p.Response
		res, err = client.Send("dialog", req)
		if err != nil {
			log.Panicln(err)
		}

		var buy Buy
		err = res.GetGob(&buy)
		if err != nil {
			log.Panicln(err)
		}

		fmt.Printf("> Buy: %s\n", buy.Text)
	}
}
```

## Running

If you run the server and the client separately then you see:

* in the server stdout:

```text
> Hello: User #1
dialog: addr (127.0.0.1:52539), read (335 µs), handle (124 µs), write (95 µs), total (555 µs)
> Hello: User #2
dialog: addr (127.0.0.1:52540), read (530 µs), handle (694 µs), write (169 µs), total (1 ms)
> Hello: User #3
dialog: addr (127.0.0.1:52541), read (410 µs), handle (194 µs), write (119 µs), total (724 µs)
> Hello: User #4
dialog: addr (127.0.0.1:52542), read (280 µs), handle (113 µs), write (51 µs), total (446 µs)
> Hello: User #5
dialog: addr (127.0.0.1:52543), read (218 µs), handle (90 µs), write (41 µs), total (350 µs)
> Hello: User #6
dialog: addr (127.0.0.1:52544), read (133 µs), handle (105 µs), write (62 µs), total (301 µs)
> Hello: User #7
dialog: addr (127.0.0.1:52545), read (267 µs), handle (78 µs), write (55 µs), total (401 µs)
> Hello: User #8
dialog: addr (127.0.0.1:52546), read (155 µs), handle (77 µs), write (40 µs), total (273 µs)
> Hello: User #9
dialog: addr (127.0.0.1:52547), read (275 µs), handle (143 µs), write (58 µs), total (477 µs)
> Hello: User #10
dialog: addr (127.0.0.1:52548), read (379 µs), handle (202 µs), write (77 µs), total (658 µs)
```

* in the client stdout:

```text
dialog: addr (127.0.0.1:8080), write (423 µs), read (1 ms), total (1 ms)
> Buy: User #1
dialog: addr (127.0.0.1:8080), write (127 µs), read (408 µs), total (536 µs)
> Buy: User #2
dialog: addr (127.0.0.1:8080), write (105 µs), read (520 µs), total (625 µs)
> Buy: User #3
dialog: addr (127.0.0.1:8080), write (43 µs), read (306 µs), total (349 µs)
> Buy: User #4
dialog: addr (127.0.0.1:8080), write (43 µs), read (382 µs), total (426 µs)
> Buy: User #5
dialog: addr (127.0.0.1:8080), write (46 µs), read (334 µs), total (380 µs)
> Buy: User #6
dialog: addr (127.0.0.1:8080), write (59 µs), read (353 µs), total (413 µs)
> Buy: User #7
dialog: addr (127.0.0.1:8080), write (71 µs), read (306 µs), total (377 µs)
> Buy: User #8
dialog: addr (127.0.0.1:8080), write (39 µs), read (329 µs), total (368 µs)
> Buy: User #9
dialog: addr (127.0.0.1:8080), write (123 µs), read (387 µs), total (510 µs)
> Buy: User #10
```

* logging

All lines that start from counter (it is a topic for the communication) are logging in StdOut.

If you want to reassign this logger you need to implement your own with the following interface:

```go
type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}
```

and set it up in your server or client implementation this way:

```go
settings.SetLogger(yourLogger)
```

## List all methods

### TCP Initialization

* p2p.NewTCP(host, port) (tcp TCP) - creates TCP connection

### Server settings initialization
* p2p.NewServerSettings() (stg) - creates a new server's settings
* stg.SetLogger(l) - reassigns server's logger
* stg.SetConnTimeout(dur) - sets connection timout
* stg.SetHandleTimeout(dur) - sets handle timout
* stg.SetBodyLimit(limit) - sets max body size for reading 

### Server
* p2p.NewServer(tcp, stg) - creates a new server
* srv.SetContext(ctx) - sets context
* srv.SetHandle(topic, handler) - sets a handler that processes all request with defined topic
* srv.Serve() (err) - starts to serve

### Client settings initialization
* p2p.NewClientSettings() (stg) - creates a new server's settings
* stg.SetLogger(l) - reassigns server's logger
* stg.SetConnTimeout(dur) - sets connection timout
* stg.SetBodyLimit(limit) - sets max body size for writing
* stg.SetRetry(retries, delay) - sets retry parameters

### Client
* NewClient(tcp, stg) (clt) - creates a new client
* clt.Send(topic, req) (res, err) - sends a request to a server by the topic

### Request
* req.SetBytes(bs) - sets bytes to the request
* req.GetBytes() (bs) - gets bytes from the request
* req.SetGob(obj) (err) - encodes to Gob and sets structure to the request
* req.GetGob(obj) (err) - decode from Gob and gets structure from the request
* req.SetJson(obj) (err) - encodes to Json and sets structure to the request
* req.GetJson(obj) (err) - decode from Json and gets structure from the request
* req.String() (str) - returns string from the request

### Response
* res.SetBytes(bs) - sets bytes to the response
* res.GetBytes() (bs) - gets bytes from the response
* res.SetGob(obj) (err) - encodes to Gob and sets structure to the response
* res.GetGob(obj) (err) - decode from Gob and gets structure from the response
* res.SetJson(obj) (err) - encodes to Json and sets structure to the response
* res.GetJson(obj) (err) - decode from Json and gets structure from the response
* res.String() (str) - returns string from the response