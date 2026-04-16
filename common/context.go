package common

import (
	"context"
	"time"
)

const DefaultRPCTimeout = 5 * time.Second

// RPCContext returns a context with a default 5-second timeout for RPC calls.
// Always defer the returned cancel function.
func RPCContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultRPCTimeout)
}
