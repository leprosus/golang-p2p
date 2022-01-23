package p2p

import "time"

type ClientSettings struct {
	Limiter
	Retry
}

func NewClientSettings() (stg *ClientSettings) {
	return &ClientSettings{
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
