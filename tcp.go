package p2p

import (
	"net"
)

type TCP struct {
	addr      string
	cipherKey *CipherKey
}

func NewTCP(host, port string) (tcp *TCP) {
	return &TCP{
		addr: net.JoinHostPort(host, port),
	}
}
