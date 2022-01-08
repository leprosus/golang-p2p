package p2p

import "time"

type ClientSettings struct {
	Limiter
	Retry
	Logger
}

func NewClientSettings() (stg ClientSettings) {
	return ClientSettings{
		Logger: NewStdLogger(),
		Limiter: Limiter{
			Timeout: Timeout{
				conn: DefaultConnTimeout,
			},
			body: DefaultBodyLimit,
		},
		Retry: Retry{
			retries: DefaultRetries,
			delay:   DefaultDelayTimeout,
		},
	}
}

func (stg *ClientSettings) SetLogger(l Logger) {
	stg.Logger = l
}

func (stg *ClientSettings) SetConnTimeout(dur time.Duration) {
	stg.Limiter.conn = dur
}

func (stg *ClientSettings) SetBodyLimit(limit uint) {
	stg.Limiter.body = int(limit)
}

func (stg *ClientSettings) SetRetry(retries uint, delay time.Duration) {
	stg.Retry.retries = retries
	stg.Retry.delay = delay
}
