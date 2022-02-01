package p2p

import (
	"context"
	"net"
	"sync"
)

type Server struct {
	tcp *TCP
	rsa *RSA

	settings *ServerSettings
	logger   Logger

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

	s.settings = NewServerSettings()

	s.rsa, err = NewRSA()

	return
}

func (s *Server) SetHandler(topic string, handler Handler) {
	s.mx.Lock()
	s.handlers[topic] = handler
	s.mx.Unlock()
}

func (s *Server) SetContext(ctx context.Context) {
	s.mx.Lock()
	s.ctx = ctx
	s.mx.Unlock()
}

func (s *Server) GetContext() (ctx context.Context) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.ctx
}

func (s *Server) SetSettings(settings *ServerSettings) {
	s.settings = settings
}

func (s *Server) SetLogger(logger Logger) {
	s.logger = logger
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

		wrapped, err = NewConn(conn, s.settings.Limiter)
		if err != nil {
			s.logger.Error(err.Error())

			return
		}

		go s.processConn(wrapped, *s.settings)
	}
}

func (s *Server) processConn(conn Conn, settings ServerSettings) {
	defer func() {
		err := conn.Close()
		if err != nil {
			s.logger.Error(err.Error())
		}
	}()

	var (
		p Package

		metrics = newMetrics(conn.RemoteAddr().String())

		err error
	)
	for {
		err = conn.ReadPackage(&p)
		if err != nil {
			s.logger.Error(err.Error())

			return
		}

		err = s.processPackage(conn, settings, p, metrics)
		if p.Type == Exchange {
			break
		}

		if err != nil {
			break
		}
	}

	if err != nil {
		s.logger.Error(err.Error())

		return
	}

	s.logger.Info(metrics.string())
}

func (s *Server) processPackage(conn Conn, settings ServerSettings, p Package, metrics *Metrics) (err error) {
	switch p.Type {
	case Handshake:
		err = s.doHandshake(conn, p, metrics)
	case Exchange:
		err = s.doExchange(conn, p, settings, metrics)
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

func (s *Server) doExchange(conn Conn, p Package, settings ServerSettings, metrics *Metrics) (err error) {
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

	var (
		ctx    context.Context
		cancel context.CancelFunc

		handler Handler
		ok      bool
	)

	s.mx.RLock()
	ctx, cancel = context.WithTimeout(s.ctx, settings.Timeout.handle)
	defer cancel()

	handler, ok = s.handlers[msg.Topic]
	s.mx.RUnlock()

	if !ok {
		s.logger.Warn(UnsupportedTopic.Error())

		return
	}

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
