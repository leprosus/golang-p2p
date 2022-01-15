package p2p

import (
	"bufio"
	"context"
	"encoding/gob"
	"net"
	"sync"
	"time"
)

type Server struct {
	tcp *TCP
	rsa *RSA
	stg *ServerSettings

	mx       sync.RWMutex
	handlers map[string]Handler

	ctx context.Context
}

func NewServer(tcp *TCP, stg *ServerSettings) (s *Server, err error) {
	s = &Server{
		tcp: tcp,
		stg: stg,

		mx:       sync.RWMutex{},
		handlers: map[string]Handler{},

		ctx: context.Background(),
	}

	s.rsa, err = NewRSA()

	return
}

func (s *Server) SetContext(ctx context.Context) {
	s.mx.Lock()
	s.ctx = ctx
	s.mx.Unlock()
}

func (s *Server) SetHandle(topic string, handler Handler) {
	s.mx.Lock()
	s.handlers[topic] = handler
	s.mx.Unlock()
}

func (s *Server) Serve() (err error) {
	var listener net.Listener
	listener, err = net.Listen("tcp", s.tcp.addr)
	if err != nil {
		return
	}

	defer func() {
		err := listener.Close()
		if err != nil {
			s.stg.Logger.Error(err.Error())
		}
	}()

	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			s.stg.Logger.Error(err.Error())

			return
		}

		go s.processConn(conn, *s.stg)
	}
}

func (s *Server) processConn(conn net.Conn, stg ServerSettings) {
	var err error

	defer func() {
		err = conn.Close()
		if err != nil {
			stg.Logger.Error(err.Error())
		}
	}()

	err = conn.SetDeadline(time.Now().Add(s.stg.Timeout.conn))
	if err != nil {
		s.stg.Logger.Warn(err.Error())

		return
	}

	metrics := newMetrics(conn.RemoteAddr().String())

	var ck CipherKey
	ck, err = s.doHandshake(conn)
	if err != nil {
		s.stg.Logger.Error(err.Error())

		return
	}

	metrics.fixHandshake()

	var cm CryptMessage
	err = gob.NewDecoder(bufio.NewReaderSize(conn, stg.Limiter.body)).Decode(&cm)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	var msg Message
	msg, err = cm.Decode(ck)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	metrics.setTopic(msg.Topic)
	metrics.fixReadDuration()

	s.mx.RLock()
	ctx := s.ctx
	handler, ok := s.handlers[msg.Topic]
	s.mx.RUnlock()
	if !ok {
		stg.Logger.Warn(UnsupportedTopic.Error())

		return
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, stg.Timeout.handle)
	defer cancel()

	var (
		req = Request{bs: msg.Content}
		res Response
	)
	res, err = handler(ctx, req)
	if err != nil {
		stg.Logger.Error(err.Error())
	}

	msg.Content = res.bs
	msg.Error = err

	metrics.fixHandleDuration()

	cm, err = msg.Encode(ck)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	err = gob.NewEncoder(conn).Encode(cm)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	metrics.fixWriteDuration()

	stg.Logger.Info(metrics.string())
}

func (s *Server) doHandshake(conn net.Conn) (ck CipherKey, err error) {
	var pk PublicKey
	err = gob.NewDecoder(conn).Decode(&pk)
	if err != nil {
		return
	}

	ck, err = NewCipherKey()
	if err != nil {
		return
	}

	var cck CryptCipherKey
	cck, err = pk.Encode(ck)
	if err != nil {
		return
	}

	err = gob.NewEncoder(conn).Encode(cck)

	return
}
