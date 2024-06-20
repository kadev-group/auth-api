package cache

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"context"
	"fmt"
	"github.com/doxanocap/pkg/errs"
	"time"
)

type RequestSessionRepository struct {
	cache    interfaces.ICacheProcessor
	provider models.AuthProvider
	ttl      time.Duration
}

func InitRequestSessionRepository(cache interfaces.ICacheProcessor) *RequestSessionRepository {
	return &RequestSessionRepository{
		cache: cache,
		ttl:   consts.RequestSessionTTl,
	}
}

func (c *RequestSessionRepository) Get(ctx context.Context, key string) (string, error) {
	key = c.constructKey(key)
	raw, err := c.cache.Get(ctx, key)
	if err != nil {
		return "", errs.Wrap("auth_seance_cache.Get", err)
	}

	return string(raw), nil
}

func (c *RequestSessionRepository) Set(ctx context.Context, key string, value string) error {
	key = c.constructKey(key)
	err := c.cache.SetWithTTL(ctx, key, []byte(value), c.ttl)
	if err != nil {
		return errs.Wrap("auth_seance_cache.Set", err)
	}
	return nil
}

func (c *RequestSessionRepository) Delete(ctx context.Context, code string) error {
	key := c.constructKey(code)
	err := c.cache.Delete(ctx, key)
	if err != nil {
		return errs.Wrap("auth_seance_cache.Delete", err)
	}
	return nil
}

func (c *RequestSessionRepository) constructKey(code string) string {
	return fmt.Sprintf("%s.%s", consts.CacheRequestSessionsPrefix, code)
}
