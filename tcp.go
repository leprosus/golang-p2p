package p2p

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	DefaultTimeout          = 30 * time.Second
	DefaultRequestLimit int = 1024
)

type TCP struct {
	mx *sync.Mutex

	addr string

	listener net.Listener
	conn     net.Conn

	timeout time.Duration
	limit   int
}

func NewTCP(host string, port uint) (tcp *TCP) {
	tcp = &TCP{
		mx:   &sync.Mutex{},
		addr: fmt.Sprintf("%s:%d", host, port),
	}

	tcp.SetTimeout(DefaultTimeout)
	tcp.SetRequestLimit(DefaultRequestLimit)

	return
}

func (tcp *TCP) SetTimeout(timeout time.Duration) {
	tcp.mx.Lock()
	defer tcp.mx.Unlock()

	tcp.timeout = timeout
}

func (tcp *TCP) GetTimeout() (timeout time.Duration) {
	tcp.mx.Lock()
	defer tcp.mx.Unlock()

	return tcp.timeout
}

func (tcp *TCP) SetRequestLimit(limit int) {
	tcp.mx.Lock()
	defer tcp.mx.Unlock()

	tcp.limit = limit
}

func (tcp *TCP) GetRequestLimit() (limit int) {
	tcp.mx.Lock()
	defer tcp.mx.Unlock()

	return tcp.limit
}

func (tcp *TCP) Close() (err error) {
	tcp.mx.Lock()
	defer tcp.mx.Unlock()

	if tcp.listener != nil {
		err = tcp.listener.Close()
		if err != nil {
			return
		}
	}

	if tcp.conn != nil {
		err = tcp.conn.Close()
		if err != nil {
			return
		}
	}

	return
}
