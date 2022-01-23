package p2p

import (
	"fmt"
	"net"
)

type TCP struct {
	addr      string
	cipherKey *CipherKey
}

func NewTCP(host string, port uint) (tcp *TCP) {
	return &TCP{
		addr: net.JoinHostPort(host, fmt.Sprint(port)),
	}
}
