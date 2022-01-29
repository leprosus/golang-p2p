package p2p

import (
	"context"
	"net"
	"sync"
)

type Server struct {
	tcp    *TCP
	rsa    *RSA
	stg    *ServerSettings
	logger Logger

	ctx context.Context

	mx       sync.RWMutex
	handlers map[string]Handler
}

func NewServer(tcp *TCP) (s *Server, err error) {
	var ck CipherKey
	ck, err = NewCipherKey()
	if err != nil {
		return
	}

	tcp.cipherKey = &ck

	s = &Server{
		tcp:    tcp,
		logger: NewStdLogger(),

		ctx: context.Background(),

		mx:       sync.RWMutex{},
		handlers: map[string]Handler{},
	}

	s.stg = NewServerSettings()

	s.rsa, err = NewRSA()

	return
}

func (s *Server) SetHandle(topic string, handler Handler) {
	s.mx.Lock()
	s.handlers[topic] = handler
	s.mx.Unlock()
}

func (s *Server) SetContext(ctx context.Context) {
	s.mx.Lock()
	s.ctx = ctx
	s.mx.Unlock()
}

func (s *Server) SetSettings(stg *ServerSettings) {
	s.mx.Lock()
	s.stg = stg
	s.mx.Unlock()
}

func (s *Server) SetLogger(logger Logger) {
	s.mx.Lock()
	s.logger = logger
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
			s.logger.Error(err.Error())
		}
	}()

	var (
		conn    net.Conn
		wrapped Conn
	)
	for {
		conn, err = listener.Accept()
		if err != nil {
			s.logger.Error(err.Error())

			return
		}

		wrapped, err = NewConn(conn, s.stg.Limiter)
		if err != nil {
			s.logger.Error(err.Error())

			return
		}

		go s.processConn(wrapped, *s.stg)
	}
}

func (s *Server) processConn(conn Conn, stg ServerSettings) {
	defer func() {
		err := conn.Close()
		if err != nil {
			s.logger.Error(err.Error())
		}
	}()

	var (
		metrics  = newMetrics(conn.RemoteAddr().String())
		isFinish bool
		err      error
	)
	for {
		isFinish, err = s.processPackage(conn, stg, metrics)
		if isFinish || err != nil {
			break
		}
	}

	if err != nil {
		s.logger.Error(err.Error())

		return
	}

	s.logger.Info(metrics.string())
}

func (s *Server) processPackage(conn Conn, stg ServerSettings, metrics *Metrics) (isFinish bool, err error) {
	var p Package
	err = conn.ReadPackage(&p)
	if err != nil {
		s.logger.Error(err.Error())

		return
	}

	switch p.Type {
	case Handshake:
		err = s.doHandshake(conn, p, metrics)
	case Exchange:
		isFinish = true

		err = s.doExchange(conn, p, stg, metrics)
	default:
		err = UnsupportedPackage
	}

	return
}

func (s *Server) doHandshake(conn Conn, p Package, metrics *Metrics) (err error) {
	var pk PublicKey
	err = p.GetGob(&pk)
	if err != nil {
		s.logger.Error(err.Error())

		return
	}

	var cck CryptCipherKey
	cck, err = pk.Encode(*s.tcp.cipherKey)
	if err != nil {
		s.logger.Error(err.Error())

		return
	}

	err = p.SetGob(cck)
	if err != nil {
		s.logger.Error(err.Error())

		return
	}

	err = conn.WritePackage(p)
	if err != nil {
		s.logger.Error(err.Error())

		return
	}

	metrics.fixHandshake()

	return
}

func (s *Server) doExchange(conn Conn, p Package, stg ServerSettings, metrics *Metrics) (err error) {
	var cm CryptMessage
	err = p.GetGob(&cm)
	if err != nil {
		s.logger.Error(err.Error())

		return
	}

	var msg Message
	msg, err = cm.Decode(*s.tcp.cipherKey)
	if err != nil {
		s.logger.Warn(err.Error())

		err = s.sendError(conn, metrics)
		if err != nil {
			s.logger.Error(err.Error())
		}

		return
	}

	metrics.setTopic(msg.Topic)
	metrics.fixReadDuration()

	s.mx.RLock()
	ctx := s.ctx
	handler, ok := s.handlers[msg.Topic]
	s.mx.RUnlock()
	if !ok {
		s.logger.Warn(UnsupportedTopic.Error())

		return
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, stg.Timeout.handle)
	defer cancel()

	var req, res Data
	req.SetBytes(msg.Content)
	res, err = handler(ctx, req)
	if err != nil {
		s.logger.Error(err.Error())
	}

	msg.Content = res.GetBytes()
	msg.Error = err

	metrics.fixHandleDuration()

	cm, err = msg.Encode(*s.tcp.cipherKey)
	if err != nil {
		s.logger.Error(err.Error())

		return
	}

	err = p.SetGob(cm)
	if err != nil {
		s.logger.Error(err.Error())

		return
	}

	err = conn.WritePackage(p)
	if err != nil {
		s.logger.Error(err.Error())

		return
	}

	metrics.fixWriteDuration()

	return
}

func (s *Server) sendError(conn Conn, metrics *Metrics) (err error) {
	p := Package{
		Type: Error,
	}

	err = conn.WritePackage(p)
	if err != nil {
		s.logger.Error(err.Error())

		return
	}

	metrics.fixWriteDuration()

	return
}
