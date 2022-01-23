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
