package request

import (
	"context"
)

type Client interface {
	Send(ctx context.Context, url string) ([]byte, error)
}
