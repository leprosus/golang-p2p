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
