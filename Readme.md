# Golang simple TCP client and server

Golang-p2p is a small client and server to make p2p communication over TCP with RSA encryption.

Main aim the package is to create an easy way of microservices communication.

## Features

| Feature                     | Description                                                                                                                                |
|-----------------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| Gob, Json and Bytes support | You can send you structure or data in binary presentation or binary serialized                                                             |
| RSA  encryption             | Every communication between a client and a server start with RSA public keys handshake.<br/>All sending data are encrypted before sending. |

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

	rsa, err := p2p.NewRSA()
	if err != nil {
		log.Panicln(err)
	}

	settings := p2p.NewServerSettings()

	server := p2p.NewServer(tcp, rsa, settings)

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

	rsa, err := p2p.NewRSA()
	if err != nil {
		log.Panicln(err)
	}

	settings := p2p.NewClientSettings()

	client := p2p.NewClient(tcp, rsa, settings)

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
dialog: addr (127.0.0.1:57573), handshake (740 µs), read (4 ms), handle (104 µs), write (181 µs), total (5 ms)
> Hello: User #2
dialog: addr (127.0.0.1:57574), handshake (267 µs), read (2 ms), handle (64 µs), write (153 µs), total (3 ms)
> Hello: User #3
dialog: addr (127.0.0.1:57575), handshake (212 µs), read (2 ms), handle (50 µs), write (106 µs), total (3 ms)
> Hello: User #4
dialog: addr (127.0.0.1:57576), handshake (215 µs), read (2 ms), handle (58 µs), write (132 µs), total (3 ms)
> Hello: User #5
dialog: addr (127.0.0.1:57577), handshake (192 µs), read (2 ms), handle (55 µs), write (120 µs), total (3 ms)
> Hello: User #6
dialog: addr (127.0.0.1:57578), handshake (295 µs), read (2 ms), handle (101 µs), write (144 µs), total (3 ms)
> Hello: User #7
dialog: addr (127.0.0.1:57579), handshake (256 µs), read (2 ms), handle (193 µs), write (273 µs), total (3 ms)
> Hello: User #8
dialog: addr (127.0.0.1:57580), handshake (396 µs), read (2 ms), handle (75 µs), write (176 µs), total (3 ms)
> Hello: User #9
dialog: addr (127.0.0.1:57581), handshake (423 µs), read (2 ms), handle (77 µs), write (148 µs), total (3 ms)
> Hello: User #10
dialog: addr (127.0.0.1:57582), handshake (335 µs), read (3 ms), handle (202 µs), write (296 µs), total (3 ms)
```

* in the client stdout:

```text
> Buy: User #1
dialog: addr (127.0.0.1:8080), handshake (294 µs), write (141 µs), read (5 ms), total (5 ms)
> Buy: User #2
dialog: addr (127.0.0.1:8080), handshake (319 µs), write (138 µs), read (5 ms), total (5 ms)
> Buy: User #3
dialog: addr (127.0.0.1:8080), handshake (301 µs), write (119 µs), read (5 ms), total (5 ms)
> Buy: User #4
dialog: addr (127.0.0.1:8080), handshake (289 µs), write (115 µs), read (5 ms), total (6 ms)
> Buy: User #5
dialog: addr (127.0.0.1:8080), handshake (313 µs), write (125 µs), read (5 ms), total (5 ms)
> Buy: User #6
dialog: addr (127.0.0.1:8080), handshake (673 µs), write (158 µs), read (5 ms), total (6 ms)
> Buy: User #7
dialog: addr (127.0.0.1:8080), handshake (313 µs), write (129 µs), read (5 ms), total (6 ms)
> Buy: User #8
dialog: addr (127.0.0.1:8080), handshake (686 µs), write (229 µs), read (5 ms), total (6 ms)
> Buy: User #9
dialog: addr (127.0.0.1:8080), handshake (405 µs), write (122 µs), read (5 ms), total (6 ms)
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

### RSA Initialization

* p2p.NewRSA() (rsa, err) - creates RSA private/public keys

### Server settings initialization

* p2p.NewServerSettings() (stg) - creates a new server's settings
* stg.SetLogger(l) - reassigns server's logger
* stg.SetConnTimeout(dur) - sets connection timout
* stg.SetHandleTimeout(dur) - sets handle timout
* stg.SetBodyLimit(limit) - sets max body size for reading

### Server

* p2p.NewServer(tcp, rsa, stg) - creates a new server
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

* NewClient(tcp, rsa, stg) (clt) - creates a new client
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