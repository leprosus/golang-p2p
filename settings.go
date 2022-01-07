package p2p

import "time"

const DefaultBodyLimit = 1024

type Limiter struct {
	Timeout
	body int
}

const (
	DefaultConnTimeout   = 250 * time.Millisecond
	DefaultHandleTimeout = 250 * time.Millisecond
)

type Timeout struct {
	conn   time.Duration
	handle time.Duration
}

const (
	DefaultRetries      = 3
	DefaultDelayTimeout = 50 * time.Millisecond
)

type Retry struct {
	retries uint
	delay   time.Duration
}
