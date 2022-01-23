package p2p

import "context"

type HandlerType uint

type Handler func(ctx context.Context, req Binary) (res Binary, err error)
