package transport

import "context"

type Transport interface {
	Serve(port int) error
	GracefulStop(ctx context.Context)
}
