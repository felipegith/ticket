package ports

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (value string, found bool, err error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	GetByField(ctx context.Context, key, field string) (string, error)
}
