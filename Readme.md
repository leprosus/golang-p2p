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

	server, err := p2p.NewServer(tcp)
	if err != nil {
		log.Panicln(err)
	}

	server.SetHandler("dialog", func(ctx context.Context, req p2p.Data) (res p2p.Data, err error) {
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

	client, err := p2p.NewClient(tcp)
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
dialog: addr (127.0.0.1:50099), handshake (2 ms), read (7 ms), handle (278 µs), write (301 µs), total (10 ms)
> Hello: User #2
dialog: addr (127.0.0.1:50100), read (743 µs), handle (238 µs), write (219 µs), total (1 ms)
> Hello: User #3
dialog: addr (127.0.0.1:50101), read (834 µs), handle (233 µs), write (228 µs), total (1 ms)
> Hello: User #4
dialog: addr (127.0.0.1:50102), read (547 µs), handle (227 µs), write (260 µs), total (1 ms)
> Hello: User #5
dialog: addr (127.0.0.1:50103), read (625 µs), handle (230 µs), write (271 µs), total (1 ms)
> Hello: User #6
dialog: addr (127.0.0.1:50104), read (602 µs), handle (241 µs), write (234 µs), total (1 ms)
> Hello: User #7
dialog: addr (127.0.0.1:50105), read (589 µs), handle (258 µs), write (227 µs), total (1 ms)
> Hello: User #8
dialog: addr (127.0.0.1:50106), read (635 µs), handle (232 µs), write (221 µs), total (1 ms)
> Hello: User #9
dialog: addr (127.0.0.1:50107), read (1 ms), handle (376 µs), write (365 µs), total (1 ms)
> Hello: User #10
dialog: addr (127.0.0.1:50108), read (635 µs), handle (370 µs), write (434 µs), total (1 ms)

```

* in the client stdout:

```text
dialog: addr (127.0.0.1:8080), handshake (8 ms), write (480 µs), read (1 ms), total (10 ms)
> Buy: User #1
dialog: addr (127.0.0.1:8080), write (342 µs), read (1 ms), total (1 ms)
> Buy: User #2
dialog: addr (127.0.0.1:8080), write (451 µs), read (1 ms), total (1 ms)
> Buy: User #3
dialog: addr (127.0.0.1:8080), write (226 µs), read (1 ms), total (1 ms)
> Buy: User #4
dialog: addr (127.0.0.1:8080), write (246 µs), read (1 ms), total (1 ms)
> Buy: User #5
dialog: addr (127.0.0.1:8080), write (262 µs), read (1 ms), total (1 ms)
> Buy: User #6
dialog: addr (127.0.0.1:8080), write (262 µs), read (1 ms), total (1 ms)
> Buy: User #7
dialog: addr (127.0.0.1:8080), write (247 µs), read (1 ms), total (1 ms)
> Buy: User #8
dialog: addr (127.0.0.1:8080), write (599 µs), read (2 ms), total (2 ms)
> Buy: User #9
dialog: addr (127.0.0.1:8080), write (259 µs), read (2 ms), total (2 ms)
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

* p2p.NewServerSettings() (settings) - creates a new server's settings
* settings.SetConnTimeout(duration) - sets connection timout
* settings.SetHandleTimeout(duration) - sets handle timout
* settings.SetBodyLimit(limit) - sets max body size for reading

### Server

* p2p.NewServer(tcp) (server, error) - creates a new server
* server.SetSettings(settings) - sets server settings
* server.SetLogger(logger) - reassigns server's logger
* server.SetHandler(topic, handler) - sets a handler that processes all request with defined topic
* server.SetContext(context) - sets context
* server.GetContext() (context) - returns context
* server.Serve() (error) - starts to serve

### Client settings initialization

* p2p.NewClientSettings() (settings) - creates a new server's settings
* settings.SetConnTimeout(duration) - sets connection timout
* settings.SetBodyLimit(limit) - sets max body size for writing
* settings.SetRetry(retries, delay) - sets retry parameters

### Client

* NewClient(tcp, settings) (client, error) - creates a new client
* client.SetSettings(settings) - sets client settings
* client.SetLogger(logger) - reassigns client's logger
* client.Send(topic, request) (response, error) - sends a request to a server by the topic

### Request and Response

* data.SetBytes(bytes) - sets bytes to the request/response
* data.GetBytes() (bytes) - gets bytes from the request/response
* data.SetGob(obj) (error) - encodes to Gob and sets structure to the request/response
* data.GetGob(obj) (error) - decode from Gob and gets structure from the request/response
* data.SetJson(obj) (error) - encodes to JSON and sets structure to the request/response
* data.GetJson(obj) (error) - decode from JSON and gets structure from the request/response
* data.String() (string) - returns string from the request/response