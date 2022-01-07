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

func main() {
	tcp := p2p.NewTCP("localhost", 8080)
	settings := p2p.NewServerSettings()

	server := p2p.NewServer(tcp, settings)
	defer func() {
		err := server.Close()
		if err != nil {
			log.Panicln(err)
		}
	}()

	var c uint
	server.Handle("counter", func(ctx context.Context, req []byte) (res []byte, err error) {
		fmt.Println(">", string(req))

		c++
		res = []byte(fmt.Sprintf("buy %d", c))

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

func main() {
	tcp := p2p.NewTCP("localhost", 8080)
	settings := p2p.NewClientSettings()

	client := p2p.NewClient(tcp, settings)

	for i := 0; i < 10; i++ {
		res, err := client.Send("counter", []byte(fmt.Sprintf("hello %d", i+1)))
		if err != nil {
			log.Panicln(err)
		}

		fmt.Println("<", string(res))
	}
}

```

## Running

If you run the server and the client separately then you see:

* in the server stdout:

```text
> hello 1
counter: addr (127.0.0.1:54014), read (344 µs), handle (26 µs), write (97 µs), total (468 µs)
> hello 2
counter: addr (127.0.0.1:54015), read (263 µs), handle (19 µs), write (50 µs), total (333 µs)
> hello 3
counter: addr (127.0.0.1:54016), read (310 µs), handle (44 µs), write (152 µs), total (507 µs)
> hello 4
counter: addr (127.0.0.1:54017), read (143 µs), handle (12 µs), write (48 µs), total (204 µs)
> hello 5
counter: addr (127.0.0.1:54018), read (118 µs), handle (12 µs), write (38 µs), total (169 µs)
> hello 6
counter: addr (127.0.0.1:54019), read (154 µs), handle (16 µs), write (51 µs), total (222 µs)
> hello 7
counter: addr (127.0.0.1:54020), read (133 µs), handle (9 µs), write (47 µs), total (190 µs)
> hello 8
counter: addr (127.0.0.1:54021), read (167 µs), handle (18 µs), write (48 µs), total (234 µs)
> hello 9
counter: addr (127.0.0.1:54022), read (212 µs), handle (20 µs), write (69 µs), total (302 µs)
> hello 10
counter: addr (127.0.0.1:54023), read (226 µs), handle (25 µs), write (47 µs), total (299 µs)
```

* in the client stdout:

```text
counter: addr (127.0.0.1:8080), write (216 µs), read (470 µs), total (687 µs)
< buy 1
counter: addr (127.0.0.1:8080), write (108 µs), read (321 µs), total (430 µs)
< buy 2
counter: addr (127.0.0.1:8080), write (151 µs), read (699 µs), total (851 µs)
< buy 3
counter: addr (127.0.0.1:8080), write (61 µs), read (295 µs), total (356 µs)
< buy 4
counter: addr (127.0.0.1:8080), write (56 µs), read (246 µs), total (303 µs)
< buy 5
counter: addr (127.0.0.1:8080), write (58 µs), read (300 µs), total (358 µs)
< buy 6
counter: addr (127.0.0.1:8080), write (55 µs), read (267 µs), total (322 µs)
< buy 7
counter: addr (127.0.0.1:8080), write (59 µs), read (327 µs), total (387 µs)
< buy 8
counter: addr (127.0.0.1:8080), write (78 µs), read (354 µs), total (432 µs)
< buy 9
counter: addr (127.0.0.1:8080), write (97 µs), read (361 µs), total (459 µs)
< buy 10
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

### Server
* p2p.NewServer(tcp, stg) - creates a new server
* srv.Handle(topic, handler) - sets a handler that processes all request with defined topic
* srv.Serve() (err) - starts to serve
* srv.Close() (err) - stops and closes the server

### Client settings initialization
* p2p.NewClientSettings() (stg) - creates a new server's settings
* stg.SetLogger(l) - reassigns server's logger
* stg.SetConnTimeout(dur) - sets connection timout
* stg.SetRetry(retries, delay) - sets retry parameters

### Client
* NewClient(tcp, stg) (clt) - creates a new client
* clt.Send(topic, req) (res, err) - sends bytes to a server by the topic
