package p2p

import (
	"context"
	"net"
	"sync"
)

type Server struct {
	tcp *TCP
	rsa *RSA
	stg *ServerSettings

	ctx context.Context

	mx       sync.RWMutex
	handlers map[string]Handler
}

func NewServer(tcp *TCP, stg *ServerSettings) (s *Server, err error) {
	var ck CipherKey
	ck, err = NewCipherKey()
	if err != nil {
		return
	}

	tcp.cipherKey = &ck

	s = &Server{
		tcp: tcp,
		stg: stg,

		ctx: context.Background(),

		mx:       sync.RWMutex{},
		handlers: map[string]Handler{},
	}

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

	var (
		conn    net.Conn
		wrapped Conn
	)
	for {
		conn, err = listener.Accept()
		if err != nil {
			s.stg.Logger.Error(err.Error())

			return
		}

		wrapped, err = NewConn(conn, s.stg.Limiter)
		if err != nil {
			s.stg.Logger.Error(err.Error())

			return
		}

		go s.processConn(wrapped, *s.stg)
	}
}

func (s *Server) processConn(conn Conn, stg ServerSettings) {
	var err error

	defer func() {
		err := conn.Close()
		if err != nil {
			stg.Logger.Error(err.Error())
		}
	}()

	var p Package
	err = conn.ReadPackage(&p)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	metrics := newMetrics(conn.RemoteAddr().String())

	switch p.Type {
	case Handshake:
		err = s.doHandshake(conn, p, stg, metrics)
	case Resume:
		err = s.doResume(conn, p, stg, metrics)
	default:
		err = UnsupportedPackage

		stg.Logger.Error(err.Error())
	}
	if err != nil {
		return
	}

	err = conn.ReadPackage(&p)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	if p.Type != Exchange {
		err = UnexpectedPackage

		stg.Logger.Error(err.Error())

		return
	}

	err = s.doExchange(conn, p, stg, metrics)
	if err != nil {
		return
	}

	stg.Logger.Info(metrics.string())
}

func (s *Server) doHandshake(conn Conn, p Package, stg ServerSettings, metrics *Metrics) (err error) {
	var pk PublicKey
	err = p.GetGob(&pk)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	var cck CryptCipherKey
	cck, err = pk.Encode(*s.tcp.cipherKey)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	err = p.SetGob(cck)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	err = conn.WritePackage(p)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	metrics.fixHandshake()

	return
}

func (s *Server) doResume(conn Conn, p Package, stg ServerSettings, metrics *Metrics) (err error) {
	var rs = ResumeImpossible

	var bs = p.GetBytes()
	_, err = s.tcp.cipherKey.Decode(bs)
	if err == nil {
		rs = ResumePossible
	}

	err = p.SetGob(rs)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	err = conn.WritePackage(p)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	metrics.fixResume()

	return
}

func (s *Server) doExchange(conn Conn, p Package, stg ServerSettings, metrics *Metrics) (err error) {
	var cm CryptMessage
	err = p.GetGob(&cm)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	var msg Message
	msg, err = cm.Decode(*s.tcp.cipherKey)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	metrics.setTopic(msg.Topic)
	metrics.fixReadDuration()

	s.mx.RLock()
	handler, ok := s.handlers[msg.Topic]
	s.mx.RUnlock()
	if !ok {
		stg.Logger.Warn(UnsupportedTopic.Error())

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), stg.Timeout.handle)
	defer cancel()

	var req, res Data
	req.SetBytes(msg.Content)
	res, err = handler(ctx, req)
	if err != nil {
		stg.Logger.Error(err.Error())
	}

	msg.Content = res.GetBytes()
	msg.Error = err

	metrics.fixHandleDuration()

	cm, err = msg.Encode(*s.tcp.cipherKey)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	err = p.SetGob(cm)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	err = conn.WritePackage(p)
	if err != nil {
		stg.Logger.Error(err.Error())

		return
	}

	metrics.fixWriteDuration()

	return
}
