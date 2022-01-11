package p2p

import (
	"bufio"
	"encoding/gob"
	"errors"
	"net"
	"sync"
	"time"
)

type Client struct {
	tcp TCP
	stg ClientSettings

	mx sync.RWMutex
}

func NewClient(tcp TCP, stg ClientSettings) (c *Client) {
	return &Client{
		tcp: tcp,
		stg: stg,

		mx: sync.RWMutex{},
	}
}

func (c *Client) SendBytes(topic string, req []byte) (res []byte, err error) {
	var retries = c.stg.retries
	for retries > 0 {
		c.mx.RLock()
		factor := c.stg.retries - retries
		c.mx.RUnlock()
		time.Sleep(time.Duration(factor) * c.stg.delay)
		retries--

		res, err = c.try(topic, req)
		if err != nil {
			c.stg.Logger.Error(err.Error())

			continue
		}

		return
	}

	return
}

func (c *Client) SendObject(topic string, req interface{}) (res interface{}, err error) {
	var reqBs []byte
	reqBs, err = Encode(req)
	if err != nil {
		return
	}

	var retries = c.stg.retries
	for retries > 0 {
		c.mx.RLock()
		factor := c.stg.retries - retries
		c.mx.RUnlock()
		time.Sleep(time.Duration(factor) * c.stg.delay)
		retries--

		var resBs []byte
		resBs, err = c.try(topic, reqBs)
		if err != nil {
			if errors.Is(err, ConnectionError) || errors.Is(err, PresetConnectionError) {
				return
			}

			c.stg.Logger.Error(err.Error())

			continue
		}

		res, err = Decode(resBs)
		if err != nil {
			return
		}

		return
	}

	return
}

func (c *Client) try(topic string, req []byte) (res []byte, err error) {
	var conn net.Conn
	conn, err = net.Dial("tcp", c.tcp.addr)
	if err != nil {
		c.stg.Logger.Error(err.Error())

		err = ConnectionError

		return
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			c.stg.Logger.Error(err.Error())
		}
	}()

	err = conn.SetDeadline(time.Now().Add(c.stg.Timeout.conn))
	if err != nil {
		c.stg.Logger.Error(err.Error())

		err = PresetConnectionError

		return
	}

	metrics := newMetrics(conn.RemoteAddr().String())
	metrics.setTopic(topic)

	msg := Message{
		Topic:   topic,
		Content: req,
	}

	err = gob.NewEncoder(conn).
		Encode(msg)
	if err != nil {
		c.stg.Logger.Error(err.Error())

		return
	}

	metrics.fixWriteDuration()

	err = gob.NewDecoder(bufio.NewReaderSize(conn, c.stg.Limiter.body)).
		Decode(&msg)
	if err != nil {
		c.stg.Logger.Error(err.Error())

		return
	}

	metrics.fixReadDuration()

	if msg.Error != nil {
		err = msg.Error

		return
	}

	c.stg.Logger.Info(metrics.string())

	res = msg.Content

	return
}
