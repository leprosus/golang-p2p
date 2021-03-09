package p2p

import (
	"bufio"
	"net"
	"time"
)

func (tcp *TCP) Handle(handler func(req *Request, res *Response)) (err error) {
	tcp.listener, err = net.Listen("tcp", tcp.addr)
	if err != nil {
		return
	}

	var conn net.Conn

	for {
		conn, err = tcp.listener.Accept()
		if err != nil {
			return
		}

		go func(conn net.Conn) {
			defer func() {
				_ = conn.Close()
			}()

			var (
				err     error
				buf     []byte
				hasMore bool
				req     = &Request{}
				res     = &Response{conn: conn}
				limit   = tcp.GetRequestLimit()
				timeout = tcp.GetTimeout()
			)

			err = conn.SetDeadline(time.Now().Add(timeout))
			if err != nil {
				return
			}

			reader := bufio.NewReaderSize(conn, limit)
			buf, hasMore, err = reader.ReadLine()

			if err != nil {
				return
			} else if hasMore {
				err = RequestTooLarge

				return
			}

			req.Body = buf

			handler(req, res)
		}(conn)
	}
}
