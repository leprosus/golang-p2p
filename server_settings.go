package p2p

import "time"

type ServerSettings struct {
	Limiter
	Logger
}

func NewServerSettings() (stg *ServerSettings) {
	return &ServerSettings{
		Logger: NewStdLogger(),
		Limiter: Limiter{
			Timeout: Timeout{
				conn:   DefaultConnTimeout,
				handle: DefaultHandleTimeout,
			},
			body: DefaultBodyLimit,
		},
	}
}

func (stg *ServerSettings) SetLogger(l Logger) {
	stg.Logger = l
}

func (stg *ServerSettings) SetConnTimeout(dur time.Duration) {
	stg.Limiter.conn = dur
}

func (stg *ServerSettings) SetHandleTimeout(dur time.Duration) {
	stg.Limiter.handle = dur
}

func (stg *ServerSettings) SetBodyLimit(limit uint) {
	stg.Limiter.body = int(limit)
}
