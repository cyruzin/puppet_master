package domain

import (
	"context"
	"time"
)

// CacheRepository represent the cache's repostiory contract.
type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, destination interface{}) error
}
