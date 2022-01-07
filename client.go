package p2p

import (
	"bufio"
	"encoding/gob"
	"net"
	"time"
)

type Client struct {
	TCP
	ClientSettings
}

func NewClient(tcp TCP, stg ClientSettings) (c *Client) {
	return &Client{
		TCP:            tcp,
		ClientSettings: stg,
	}
}

func (c *Client) Send(topic string, req []byte) (res []byte, err error) {
	var retries = c.ClientSettings.retries
	for retries > 0 {
		factor := c.ClientSettings.retries - retries
		time.Sleep(time.Duration(factor) * c.ClientSettings.delay)
		retries--

		res, err = c.try(topic, req)
		if err != nil {
			c.ClientSettings.Logger.Error(err.Error())

			continue
		}

		return
	}

	return
}

func (c *Client) try(topic string, req []byte) (res []byte, err error) {
	var conn net.Conn
	conn, err = net.Dial("tcp", c.TCP.addr)
	if err != nil {
		return
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			c.Logger.Error(err.Error())
		}
	}()

	err = conn.SetDeadline(time.Now().Add(c.Timeout.conn))
	if err != nil {
		c.Logger.Warn(err.Error())

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
		c.Logger.Error(err.Error())

		return
	}

	metrics.fixWriteDuration()

	err = gob.NewDecoder(bufio.NewReaderSize(conn, c.Limiter.body)).
		Decode(&msg)
	if err != nil {
		c.Logger.Error(err.Error())

		return
	}

	metrics.fixReadDuration()

	c.Logger.Info(metrics.string())

	res = msg.Content

	return
}
