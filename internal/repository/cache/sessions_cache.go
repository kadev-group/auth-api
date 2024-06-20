package cache

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models/consts"
	"context"
	"fmt"
	"github.com/doxanocap/pkg/errs"
	"strconv"
	"time"
)

type SessionCacheRepository struct {
	cache interfaces.ICacheProcessor
	ttl   time.Duration
}

func InitSessionCacheRepository(cache interfaces.ICacheProcessor) *SessionCacheRepository {
	return &SessionCacheRepository{
		cache: cache,
		ttl:   consts.AccessTokenTTL,
	}
}

func (c *SessionCacheRepository) Get(ctx context.Context, userIDCode string) (int64, error) {
	key := c.constructKey(userIDCode)

	raw, err := c.cache.Get(ctx, key)
	if err != nil {
		return 0, errs.Wrap("sessions_cache.Get", err)
	}
	if raw == nil {
		return 0, nil
	}

	value, err := strconv.Atoi(string(raw))
	if err != nil {
		return 0, errs.Wrap("sessions_cache.Get: convert", err)
	}

	return int64(value), nil
}

func (c *SessionCacheRepository) Set(ctx context.Context, userIDCode string, startedAt int64) error {
	key := c.constructKey(userIDCode)

	value := strconv.Itoa(int(startedAt))
	err := c.cache.Set(ctx, key, []byte(value))
	if err != nil {
		return errs.Wrap("sessions_cache.Set", err)
	}
	return nil
}

func (c *SessionCacheRepository) constructKey(userIDCode string) string {
	return fmt.Sprintf("%s.%s", consts.CacheSessionsPrefix, userIDCode)
}
