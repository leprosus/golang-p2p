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
	tcp := p2p.NewTCP("localhost", "8080")

	settings := p2p.NewServerSettings()

	server, err := p2p.NewServer(tcp, settings)
	if err != nil {
		log.Panicln(err)
	}

	server.SetHandle("dialog", func(ctx context.Context, req p2p.Data) (res p2p.Data, err error) {
		hello := Hello{}
		err = req.GetGob(&hello)
		if err != nil {
			return
		}

		fmt.Printf("> Hello: %s\n", hello.Text)

		res = p2p.Data{}
		err = res.SetGob(Buy{
			Text: hello.Text,
		})

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
	tcp := p2p.NewTCP("localhost", "8080")

	settings := p2p.NewClientSettings()

	client, err := p2p.NewClient(tcp, settings)
	if err != nil {
		log.Panicln(err)
	}

	var req, res p2p.Data

	for i := 0; i < 10; i++ {
		req = p2p.Data{}
		err = req.SetGob(Hello{
			Text: fmt.Sprintf("User #%d", i+1),
		})
		if err != nil {
			log.Panicln(err)
		}

		res = p2p.Data{}
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
dialog: addr (127.0.0.1:54949), handshake (426 µs), read (4 ms), handle (240 µs), write (215 µs), total (5 ms)
> Hello: User #2
dialog: addr (127.0.0.1:54950), resume (85 µs), read (456 µs), handle (61 µs), write (65 µs), total (669 µs)
> Hello: User #3
dialog: addr (127.0.0.1:54951), resume (65 µs), read (409 µs), handle (71 µs), write (73 µs), total (619 µs)
> Hello: User #4
dialog: addr (127.0.0.1:54952), resume (58 µs), read (237 µs), handle (48 µs), write (56 µs), total (400 µs)
> Hello: User #5
dialog: addr (127.0.0.1:54953), resume (62 µs), read (209 µs), handle (45 µs), write (59 µs), total (377 µs)
> Hello: User #6
dialog: addr (127.0.0.1:54954), resume (51 µs), read (284 µs), handle (74 µs), write (59 µs), total (469 µs)
> Hello: User #7
dialog: addr (127.0.0.1:54955), resume (90 µs), read (352 µs), handle (97 µs), write (96 µs), total (638 µs)
> Hello: User #8
dialog: addr (127.0.0.1:54956), resume (60 µs), read (310 µs), handle (77 µs), write (102 µs), total (550 µs)
> Hello: User #9
dialog: addr (127.0.0.1:54957), resume (110 µs), read (319 µs), handle (84 µs), write (102 µs), total (617 µs)
> Hello: User #10
dialog: addr (127.0.0.1:54958), resume (75 µs), read (496 µs), handle (97 µs), write (91 µs), total (761 µs)
```

* in the client stdout:

```text
dialog: addr (127.0.0.1:8080), handshake (4 ms), write (616 µs), read (1 ms), total (6 ms)
> Buy: User #1
dialog: addr (127.0.0.1:8080), resume (620 µs), write (158 µs), read (392 µs), total (1 ms)
> Buy: User #2
dialog: addr (127.0.0.1:8080), resume (334 µs), write (56 µs), read (566 µs), total (956 µs)
> Buy: User #3
dialog: addr (127.0.0.1:8080), resume (318 µs), write (53 µs), read (329 µs), total (701 µs)
> Buy: User #4
dialog: addr (127.0.0.1:8080), resume (273 µs), write (54 µs), read (350 µs), total (678 µs)
> Buy: User #5
dialog: addr (127.0.0.1:8080), resume (331 µs), write (72 µs), read (408 µs), total (812 µs)
> Buy: User #6
dialog: addr (127.0.0.1:8080), resume (368 µs), write (157 µs), read (522 µs), total (1 ms)
> Buy: User #7
dialog: addr (127.0.0.1:8080), resume (384 µs), write (66 µs), read (480 µs), total (931 µs)
> Buy: User #8
dialog: addr (127.0.0.1:8080), resume (443 µs), write (105 µs), read (536 µs), total (1 ms)
> Buy: User #9
dialog: addr (127.0.0.1:8080), resume (539 µs), write (104 µs), read (588 µs), total (1 ms)
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
* srv.SetHandle(topic, handler) - sets a handler that processes all request with defined topic
* srv.SetContext(ctx) - sets context
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