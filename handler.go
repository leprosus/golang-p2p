package p2p

import "context"

type HandlerType uint

type Handler func(ctx context.Context, req Data) (res Data, err error)
