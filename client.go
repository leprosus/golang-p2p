package p2p

import (
	"bufio"
	"net"
)

func (tcp *TCP) Send(bytes []byte) (result []byte, err error) {
	if len(bytes) > tcp.limit {
		err = RequestTooLarge

		return
	}

	tcp.conn, err = net.Dial("tcp", tcp.addr)
	if err != nil {
		return
	}

	_, err = tcp.conn.Write(append(bytes, '\n'))
	if err != nil {
		return
	}

	reader := bufio.NewReader(tcp.conn)
	result, _, err = reader.ReadLine()
	if err != nil {
		return
	}

	err = tcp.conn.Close()

	return
}
