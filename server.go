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
	TCP
	ServerSettings

	mx       sync.RWMutex
	listener net.Listener
	handlers map[string]Handler

	ctx context.Context
}

func NewServer(tcp TCP, stg ServerSettings) (s *Server) {
	return &Server{
		TCP:            tcp,
		ServerSettings: stg,

		mx:       sync.RWMutex{},
		handlers: map[string]Handler{},

		ctx: context.Background(),
	}
}

func (s *Server) SetContext(ctx context.Context) {
	s.mx.Lock()
	s.ctx = ctx
	s.mx.Unlock()
}

func (s *Server) Handle(topic string, handler Handler) {
	s.mx.Lock()
	s.handlers[topic] = handler
	s.mx.Unlock()
}

func (s *Server) Serve() (err error) {
	var listener net.Listener
	listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return
	}

	defer func() {
		err := listener.Close()
		if err != nil {
			s.Logger.Error(err.Error())
		}
	}()

	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			s.Logger.Error(err.Error())

			return
		}

		go func(conn net.Conn, stg ServerSettings) {
			defer func() {
				err = conn.Close()
				if err != nil {
					stg.Logger.Error(err.Error())
				}
			}()

			err := conn.SetDeadline(time.Now().Add(s.Timeout.conn))
			if err != nil {
				s.Logger.Warn(err.Error())

				return
			}

			metrics := newMetrics(conn.RemoteAddr().String())

			var msg Message
			err = gob.NewDecoder(bufio.NewReaderSize(conn, stg.Limiter.body)).
				Decode(&msg)
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

			ctx, cancel := context.WithTimeout(ctx, stg.Timeout.handle)
			defer cancel()

			msg.Content, err = handler(ctx, msg.Content)
			if err != nil {
				stg.Logger.Error(err.Error())

				return
			}

			metrics.fixHandleDuration()

			err = gob.NewEncoder(conn).
				Encode(msg)
			if err != nil {
				stg.Logger.Error(err.Error())

				return
			}

			metrics.fixWriteDuration()

			stg.Logger.Info(metrics.string())
		}(conn, s.ServerSettings)
	}
}

func (s *Server) Close() (err error) {
	return s.listener.Close()
}
