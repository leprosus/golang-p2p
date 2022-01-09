package p2p

import "context"

type HandlerType uint

const (
	BytesHandlerType HandlerType = iota
	ObjectHandlerType
)

type BytesHandler func(ctx context.Context, req []byte) (res []byte, err error)

type ObjectHandler func(ctx context.Context, req interface{}) (res interface{}, err error)
