package p2p

import (
	"net"
	"sync"
	"time"
)

type Client struct {
	tcp *TCP
	rsa *RSA

	settings *ClientSettings
	logger   Logger

	mx sync.RWMutex
}

func NewClient(tcp *TCP) (c *Client, err error) {
	c = &Client{
		tcp:    tcp,
		logger: NewStdLogger(),

		mx: sync.RWMutex{},
	}

	c.settings = NewClientSettings()

	c.rsa, err = NewRSA()

	return
}

func (c *Client) SetSettings(settings *ClientSettings) {
	c.settings = settings
}

func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
}

func (c *Client) Send(topic string, req Data) (res Data, err error) {
	var retries = c.settings.retries
	for retries > 0 {
		c.mx.RLock()
		factor := c.settings.retries - retries
		c.mx.RUnlock()
		time.Sleep(time.Duration(factor) * c.settings.delay)
		retries--

		res, err = c.try(topic, req)
		if err != nil {
			continue
		}

		return
	}

	return
}

func (c *Client) try(topic string, req Data) (res Data, err error) {
	var conn net.Conn
	conn, err = net.Dial("tcp", c.tcp.addr)
	if err != nil {
		c.logger.Error(err.Error())

		err = ConnectionError

		return
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			c.logger.Error(err.Error())
		}
	}()

	var wrapped Conn
	wrapped, err = NewConn(conn, c.settings.Limiter)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	metrics := newMetrics(conn.RemoteAddr().String())
	metrics.setTopic(topic)

	msg := Message{
		Topic:   topic,
		Content: req.GetBytes(),
	}

	for {
		if c.tcp.cipherKey == nil {
			var ck CipherKey
			ck, err = c.doHandshake(wrapped, metrics)
			if err != nil {
				break
			}

			c.tcp.cipherKey = &ck
		} else {
			msg, err = c.doExchange(wrapped, metrics, msg)
			if err != nil {
				c.tcp.cipherKey = nil

				return
			}

			res.SetBytes(msg.Content)

			break
		}

		if err != nil {
			return
		}
	}

	c.logger.Info(metrics.string())

	return
}

func (c *Client) doHandshake(conn Conn, metrics *Metrics) (ck CipherKey, err error) {
	p := Package{
		Type: Handshake,
	}

	err = p.SetGob(c.rsa.PublicKey())
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	err = conn.WritePackage(p)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	err = conn.ReadPackage(&p)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	var cck CryptCipherKey
	err = p.GetGob(&cck)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	ck, err = c.rsa.PrivateKey().Decode(cck)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	metrics.fixHandshake()

	return
}

func (c *Client) doExchange(conn Conn, metrics *Metrics, in Message) (out Message, err error) {
	var cm CryptMessage
	cm, err = in.Encode(*c.tcp.cipherKey)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	p := Package{
		Type: Exchange,
	}
	err = p.SetGob(cm)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	err = conn.WritePackage(p)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	metrics.fixWriteDuration()

	err = conn.ReadPackage(&p)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	if p.Type == Error {
		_ = p.GetGob(&err)

		return
	}

	err = p.GetGob(&cm)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	out, err = cm.Decode(*c.tcp.cipherKey)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	metrics.fixReadDuration()

	if out.Error != nil {
		err = out.Error

		c.logger.Error(err.Error())

		return
	}

	return
}
