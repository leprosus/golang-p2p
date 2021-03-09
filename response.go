package p2p

import (
	"net"
)

type Response struct {
	conn net.Conn
}

func (res *Response) Send(bs []byte) (err error) {
	_, err = res.conn.Write(append(bs, '\n'))

	return
}
