package p2p

import (
	"bufio"
	"encoding/gob"
	"net"
	"time"
)

type Conn struct {
	net.Conn
	limiter Limiter
}

func NewConn(conn net.Conn, limiter Limiter) (c Conn, err error) {
	c = Conn{
		Conn:    conn,
		limiter: limiter,
	}

	err = conn.SetDeadline(time.Now().Add(limiter.Timeout.conn))
	if err != nil {
		err = PresetConnectionError

		return
	}

	return
}

func (c *Conn) ReadPackage(p *Package) (err error) {
	var reader *bufio.Reader
	if c.limiter.body > 0 {
		reader = bufio.NewReaderSize(c, c.limiter.body)
	} else {
		reader = bufio.NewReader(c)
	}

	err = gob.NewDecoder(reader).Decode(&p)

	return
}

func (c *Conn) WritePackage(p Package) (err error) {
	err = gob.NewEncoder(c).Encode(p)

	return
}
