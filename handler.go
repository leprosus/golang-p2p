package p2p

import "context"

type Handler func(ctx context.Context, req Data) (res Data, err error)
