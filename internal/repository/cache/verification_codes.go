package cache

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"context"
	"fmt"
	"github.com/doxanocap/pkg/errs"
	"github.com/go-redis/redis/v8"
	"time"
)

type VerificationCodesRepository struct {
	cache interfaces.ICacheProcessor
	ttl   time.Duration
}

func InitVerificationCodesRepository(cache interfaces.ICacheProcessor) *VerificationCodesRepository {
	return &VerificationCodesRepository{
		cache: cache,
		ttl:   consts.VerificationCodesTTL,
	}
}

func (c *VerificationCodesRepository) Get(ctx context.Context, key string) (*models.VerificationCode, error) {
	key = c.constructKey(key)

	val := &models.VerificationCode{}
	err := c.cache.GetJSON(ctx, key, val)
	if err != nil {
		if err == redis.Nil {
			return val, nil
		}
		return nil, errs.Wrap("verification_codes.Get", err)
	}
	return val, nil
}

func (c *VerificationCodesRepository) Set(ctx context.Context, key string, val *models.VerificationCode) error {
	key = c.constructKey(key)

	err := c.cache.SetJSON(ctx, key, val)
	if err != nil {
		return errs.Wrap("verification_codes.Set", err)
	}

	return nil
}

func (c *VerificationCodesRepository) Delete(ctx context.Context, key string) error {
	key = c.constructKey(key)

	err := c.cache.Delete(ctx, key)
	if err != nil {
		return errs.Wrap("verification_codes.Delete", err)
	}

	return nil
}

func (c *VerificationCodesRepository) constructKey(key string) string {
	return fmt.Sprintf("%s:%s", consts.CacheVerifyCodePrefix, key)
}
