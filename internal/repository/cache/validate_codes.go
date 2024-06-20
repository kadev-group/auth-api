package cache

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"context"
	"errors"
	"fmt"
	"github.com/doxanocap/pkg/errs"
	"github.com/go-redis/redis/v8"
	"time"
)

type ValidateCodesRepository struct {
	cache interfaces.ICacheProcessor
	ttl   time.Duration
}

func InitValidateCodesRepository(cache interfaces.ICacheProcessor) *ValidateCodesRepository {
	return &ValidateCodesRepository{
		cache: cache,
		ttl:   consts.ValidateCodesTTL,
	}
}

func (c *ValidateCodesRepository) Get(ctx context.Context, key string) (*models.VerificationCode, error) {
	key = c.constructKey(key)

	val := &models.VerificationCode{}
	err := c.cache.GetJSON(ctx, key, val)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return val, nil
		}
		return nil, errs.Wrap("verification_codes.Get", err)
	}
	return val, nil
}

func (c *ValidateCodesRepository) Set(ctx context.Context, key string, val *models.VerificationCode) error {
	key = c.constructKey(key)

	err := c.cache.SetJSONWithTTL(ctx, key, val, c.ttl)
	if err != nil {
		return errs.Wrap("verification_codes.Set", err)
	}

	return nil
}

func (c *ValidateCodesRepository) Delete(ctx context.Context, key string) error {
	key = c.constructKey(key)

	err := c.cache.Delete(ctx, key)
	if err != nil {
		return errs.Wrap("verification_codes.Delete", err)
	}

	return nil
}

func (c *ValidateCodesRepository) constructKey(key string) string {
	return fmt.Sprintf("%s.%s", consts.CacheVerifyCodePrefix, key)
}
