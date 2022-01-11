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
	tcp TCP
	stg ServerSettings

	mx       sync.RWMutex
	listener net.Listener
	handlers map[string]handler

	ctx context.Context
}

type handler struct {
	Type      HandlerType
	Interface interface{}
}

func NewServer(tcp TCP, stg ServerSettings) (s *Server) {
	return &Server{
		tcp: tcp,
		stg: stg,

		mx:       sync.RWMutex{},
		handlers: map[string]handler{},

		ctx: context.Background(),
	}
}

func (s *Server) SetContext(ctx context.Context) {
	s.mx.Lock()
	s.ctx = ctx
	s.mx.Unlock()
}

func (s *Server) SetBytesHandle(topic string, bh BytesHandler) {
	s.mx.Lock()
	s.handlers[topic] = handler{
		Type:      BytesHandlerType,
		Interface: bh,
	}
	s.mx.Unlock()
}

func (s *Server) SetObjectHandle(topic string, oh ObjectHandler) {
	s.mx.Lock()
	s.handlers[topic] = handler{
		Type:      ObjectHandlerType,
		Interface: oh,
	}
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

		go func(conn net.Conn, stg ServerSettings) {
			defer func() {
				err = conn.Close()
				if err != nil {
					stg.Logger.Error(err.Error())
				}
			}()

			err := conn.SetDeadline(time.Now().Add(s.stg.Timeout.conn))
			if err != nil {
				s.stg.Logger.Warn(err.Error())

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

			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, stg.Timeout.handle)
			defer cancel()

			var reqObj interface{}

			switch handler.Type {
			case BytesHandlerType:
				msg.Content, msg.Error = handler.Interface.(BytesHandler)(ctx, msg.Content)
				if msg.Error != nil {
					stg.Logger.Error(msg.Error.Error())

					return
				}
			case ObjectHandlerType:
				if reqObj == nil {
					reqObj, err = Decode(msg.Content)
					if err != nil {
						stg.Logger.Error(msg.Error.Error())

						return
					}
				}

				var resObj interface{}
				resObj, msg.Error = handler.Interface.(ObjectHandler)(ctx, reqObj)
				if msg.Error != nil {
					stg.Logger.Error(msg.Error.Error())

					return
				}

				msg.Content, err = Encode(resObj)
				if err != nil {
					stg.Logger.Error(msg.Error.Error())

					return
				}
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
		}(conn, s.stg)
	}
}

func (s *Server) Close() (err error) {
	return s.listener.Close()
}
