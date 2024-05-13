package cache

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"context"
	"fmt"
	"github.com/doxanocap/pkg/errs"
	"log"
	"strconv"
	"time"
)

type OAuthCacheRepository struct {
	cache    interfaces.ICacheProcessor
	provider models.OAuthProvider
	ttl      time.Duration
}

func InitOAuthCacheRepository(
	provider models.OAuthProvider, cache interfaces.ICacheProcessor) *OAuthCacheRepository {
	if !provider.IsValid() {
		log.Fatalf("invalid provider: %s", provider)
	}
	return &OAuthCacheRepository{
		cache:    cache,
		provider: provider,
		ttl:      consts.OAuthCodeTTl,
	}
}

func (c *OAuthCacheRepository) Get(ctx context.Context, code string) (int64, error) {
	key := c.constructKey(code)
	raw, err := c.cache.Get(ctx, key)
	if err != nil {
		return 0, errs.Wrap("oauth_cache.Get", err)
	}
	rawStr := string(raw)
	if rawStr == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(rawStr)
	if err != nil {
		return 0, errs.Wrap("oauth_cache.Get: convert", err)
	}
	return int64(value), nil
}

func (c *OAuthCacheRepository) Set(ctx context.Context, code string, value int64) error {
	key := c.constructKey(code)
	rawValue := strconv.Itoa(int(value))
	err := c.cache.Set(ctx, key, []byte(rawValue))
	if err != nil {
		return errs.Wrap("oauth_cache.Set", err)
	}
	return nil
}

func (c *OAuthCacheRepository) Delete(ctx context.Context, code string) error {
	key := c.constructKey(code)
	err := c.cache.Delete(ctx, key)
	if err != nil {
		return errs.Wrap("oauth_cache.Delete", err)
	}
	return nil
}

func (c *OAuthCacheRepository) constructKey(code string) string {
	switch c.provider {
	case models.GoogleOAuth:
		return fmt.Sprintf("%s:%s", consts.CacheGoogleOAuthPrefix, code)
	default:
		return ""
	}
}
