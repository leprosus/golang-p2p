# Golang simple TCP client and server

Golang-p2p is a small client and server to make p2p communication over TCP with RSA encryption.

Main aim the package is to create an easy way of microservices communication.

## Features

| Feature                     | Description                                                                                                                                 |
|-----------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|
| Gob, Json and Bytes support | You can send you structure or data in binary presentation or binary serialized                                                              |
| RSA  handshake              | Every communication between a client and a server starts with RSA public keys handshake.<br/>All sending data are encrypted before sending. |

## Import

```go
import "github.com/leprosus/golang-p2p"
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

	server, err := p2p.NewServer(tcp, settings)
	if err != nil {
		log.Panicln(err)
	}

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

	err = server.Serve()
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

	client, err := p2p.NewClient(tcp, settings)
	if err != nil {
		log.Panicln(err)
	}

	for i := 0; i < 10; i++ {
		hello := Hello{Text: fmt.Sprintf("User #%d", i+1)}

		req := p2p.Request{}
		err = req.SetGob(hello)
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
dialog: addr (127.0.0.1:61620), handshake (511 µs), read (2 ms), handle (129 µs), write (62 µs), total (3 ms)
> Hello: User #2
dialog: addr (127.0.0.1:61621), handshake (310 µs), read (2 ms), handle (60 µs), write (30 µs), total (3 ms)
> Hello: User #3
dialog: addr (127.0.0.1:61623), handshake (220 µs), read (2 ms), handle (68 µs), write (26 µs), total (3 ms)
> Hello: User #4
dialog: addr (127.0.0.1:61624), handshake (252 µs), read (2 ms), handle (79 µs), write (30 µs), total (3 ms)
> Hello: User #5
dialog: addr (127.0.0.1:61625), handshake (340 µs), read (2 ms), handle (75 µs), write (41 µs), total (3 ms)
> Hello: User #6
dialog: addr (127.0.0.1:61626), handshake (276 µs), read (2 ms), handle (57 µs), write (34 µs), total (3 ms)
> Hello: User #7
dialog: addr (127.0.0.1:61627), handshake (251 µs), read (3 ms), handle (163 µs), write (65 µs), total (3 ms)
> Hello: User #8
dialog: addr (127.0.0.1:61628), handshake (268 µs), read (2 ms), handle (89 µs), write (50 µs), total (3 ms)
> Hello: User #9
dialog: addr (127.0.0.1:61629), handshake (413 µs), read (3 ms), handle (229 µs), write (85 µs), total (4 ms)
> Hello: User #10
dialog: addr (127.0.0.1:61630), handshake (663 µs), read (3 ms), handle (88 µs), write (73 µs), total (4 ms)
```

* in the client stdout:

```text
dialog: addr (127.0.0.1:8080), handshake (3 ms), write (90 µs), read (630 µs), total (3 ms)
> Buy: User #1
dialog: addr (127.0.0.1:8080), handshake (2 ms), write (34 µs), read (476 µs), total (3 ms)
> Buy: User #2
dialog: addr (127.0.0.1:8080), handshake (2 ms), write (33 µs), read (331 µs), total (3 ms)
> Buy: User #3
dialog: addr (127.0.0.1:8080), handshake (2 ms), write (43 µs), read (369 µs), total (3 ms)
> Buy: User #4
dialog: addr (127.0.0.1:8080), handshake (2 ms), write (60 µs), read (415 µs), total (3 ms)
> Buy: User #5
dialog: addr (127.0.0.1:8080), handshake (2 ms), write (48 µs), read (374 µs), total (3 ms)
> Buy: User #6
dialog: addr (127.0.0.1:8080), handshake (2 ms), write (115 µs), read (800 µs), total (3 ms)
> Buy: User #7
dialog: addr (127.0.0.1:8080), handshake (2 ms), write (54 µs), read (560 µs), total (3 ms)
> Buy: User #8
dialog: addr (127.0.0.1:8080), handshake (3 ms), write (198 µs), read (992 µs), total (4 ms)
> Buy: User #9
dialog: addr (127.0.0.1:8080), handshake (3 ms), write (53 µs), read (502 µs), total (4 ms)
> Buy: User #10
```

* logging

All lines that start from `dialog` is the topic for the communication.

All log lines write to StdOut.

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

* p2p.NewTCP(host, port) (tcp, err) - creates TCP connection

### Server settings initialization

* p2p.NewServerSettings() (stg) - creates a new server's settings
* stg.SetLogger(l) - reassigns server's logger
* stg.SetConnTimeout(dur) - sets connection timout
* stg.SetHandleTimeout(dur) - sets handle timout
* stg.SetBodyLimit(limit) - sets max body size for reading

### Server

* p2p.NewServer(tcp, stg) (srv, err) - creates a new server
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

* NewClient(tcp, stg) (clt, err) - creates a new client
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