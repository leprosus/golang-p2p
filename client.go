package p2p

import (
	"net"
	"sync"
	"time"
)

type Client struct {
	tcp    *TCP
	rsa    *RSA
	stg    *ClientSettings
	logger Logger

	mx sync.RWMutex
}

func NewClient(tcp *TCP) (c *Client, err error) {
	c = &Client{
		tcp:    tcp,
		logger: NewStdLogger(),

		mx: sync.RWMutex{},
	}

	c.stg = NewClientSettings()

	c.rsa, err = NewRSA()

	return
}

func (c *Client) SetSettings(stg *ClientSettings) {
	c.mx.Lock()
	c.stg = stg
	c.mx.Unlock()
}

func (c *Client) SetLogger(logger Logger) {
	c.mx.Lock()
	c.logger = logger
	c.mx.Unlock()
}

func (c *Client) Send(topic string, req Data) (res Data, err error) {
	var retries = c.stg.retries
	for retries > 0 {
		c.mx.RLock()
		factor := c.stg.retries - retries
		c.mx.RUnlock()
		time.Sleep(time.Duration(factor) * c.stg.delay)
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
	wrapped, err = NewConn(conn, c.stg.Limiter)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	metrics := newMetrics(conn.RemoteAddr().String())
	metrics.setTopic(topic)

	if c.tcp.cipherKey == nil {
		var ck CipherKey
		ck, err = c.doHandshake(wrapped, metrics)
		if err != nil {
			c.logger.Error(err.Error())

			return
		}

		c.tcp.cipherKey = &ck
	} else {
		err = c.doResume(wrapped, metrics)
		if err != nil {
			c.tcp.cipherKey = nil

			return
		}
	}

	msg := Message{
		Topic:   topic,
		Content: req.GetBytes(),
	}
	msg, err = c.doExchange(wrapped, metrics, msg)
	if err != nil {
		return
	}

	res.SetBytes(msg.Content)

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

func (c *Client) doResume(conn Conn, metrics *Metrics) (err error) {
	p := Package{
		Type: Resume,
	}

	var bs []byte
	bs, err = c.tcp.cipherKey.Encode([]byte(metrics.topic))
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	p.SetBytes(bs)

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

	var rs ResumeStatus
	err = p.GetGob(&rs)
	if err != nil {
		c.logger.Error(err.Error())

		return
	}

	if rs == ResumeImpossible {
		err = CipherKeyError

		c.logger.Warn(err.Error())

		return
	}

	metrics.fixResume()

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
