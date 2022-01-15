package p2p

import "fmt"

type TCP struct {
	addr string
}

func NewTCP(host string, port uint) (tcp *TCP) {
	return &TCP{
		addr: fmt.Sprintf("%s:%d", host, port),
	}
}
