package p2p

import (
	"bufio"
	"encoding/gob"
	"net"
	"sync"
	"time"
)

type Client struct {
	tcp *TCP
	rsa *RSA
	stg *ClientSettings

	mx sync.RWMutex
}

func NewClient(tcp *TCP, stg *ClientSettings) (c *Client, err error) {
	c = &Client{
		tcp: tcp,
		stg: stg,

		mx: sync.RWMutex{},
	}

	c.rsa, err = NewRSA()

	return
}

func (c *Client) Send(topic string, req Request) (res Response, err error) {
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

func (c *Client) try(topic string, req Request) (res Response, err error) {
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

	var ck CipherKey
	ck, err = c.doHandshake(conn)
	if err != nil {
		c.stg.Logger.Error(err.Error())

		return
	}

	metrics.fixHandshake()

	msg := Message{
		Topic:   topic,
		Content: req.bs,
	}

	var cm CryptMessage
	cm, err = msg.Encode(ck)
	if err != nil {
		c.stg.Logger.Error(err.Error())

		return
	}

	err = gob.NewEncoder(conn).Encode(cm)
	if err != nil {
		c.stg.Logger.Error(err.Error())

		return
	}

	metrics.fixWriteDuration()

	err = gob.NewDecoder(bufio.NewReaderSize(conn, c.stg.Limiter.body)).Decode(&cm)
	if err != nil {
		c.stg.Logger.Error(err.Error())

		return
	}

	msg, err = cm.Decode(ck)
	if err != nil {
		c.stg.Logger.Error(err.Error())

		return
	}

	metrics.fixReadDuration()

	if msg.Error != nil {
		err = msg.Error

		c.stg.Logger.Error(err.Error())

		return
	}

	c.stg.Logger.Info(metrics.string())

	res.bs = msg.Content

	return
}

func (c *Client) doHandshake(conn net.Conn) (ck CipherKey, err error) {
	err = gob.NewEncoder(conn).Encode(c.rsa.PublicKey())
	if err != nil {
		return
	}

	var cck CryptCipherKey
	err = gob.NewDecoder(conn).Decode(&cck)
	if err != nil {
		return
	}

	ck, err = c.rsa.PrivateKey().Decode(cck)

	return
}
