package p2p

import "context"

type Handler func(ctx context.Context, req []byte) (res []byte, err error)
