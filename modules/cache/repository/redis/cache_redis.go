package rds

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type cacheRepository struct {
	Conn *redis.Client
}

// NewRedisCacheRepository will create an object that represent
// the cache.Repository interface.
func NewRedisCacheRepository(Conn *redis.Client) domain.CacheRepository {
	return &cacheRepository{Conn}
}

func (r *cacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.Conn.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Error().Err(err).Stack().Msg(domain.ErrSetCache.Error())
		return domain.ErrSetCache
	}

	return nil
}

func (r *cacheRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Conn.Get(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Stack().Msg(domain.ErrGetCache.Error())
		return "", domain.ErrGetCache
	}

	if err == redis.Nil {
		log.Error().Err(err).Stack().Msg(domain.ErrCacheKey.Error())
		return "", domain.ErrCacheKey
	}

	return val, nil
}

func (r *cacheRepository) Marshal(data interface{}) ([]byte, error) {
	value, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Stack().Msg(domain.ErrCacheMarshalling.Error())
		return nil, domain.ErrCacheMarshalling
	}

	return value, nil
}

func (r *cacheRepository) Unmarshal(data []byte, destination interface{}) error {
	err := json.Unmarshal(data, &destination)
	if err != nil {
		log.Error().Err(err).Stack().Msg(domain.ErrCacheUnmarshalling.Error())
		return domain.ErrCacheUnmarshalling
	}

	return nil
}
