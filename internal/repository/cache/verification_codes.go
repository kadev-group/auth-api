package cache

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models/consts"
	"context"
	"fmt"
	"github.com/doxanocap/pkg/errs"
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

func (c *VerificationCodesRepository) Get(ctx context.Context, email string) (string, error) {
	key := c.constructKey(email)

	raw, err := c.cache.Get(ctx, key)
	if err != nil {
		return "", errs.Wrap("verification_codes.Get", err)
	}
	return string(raw), nil
}

func (c *VerificationCodesRepository) Set(ctx context.Context, email, code string) error {
	key := c.constructKey(email)

	err := c.cache.Set(ctx, key, []byte(code))
	if err != nil {
		return errs.Wrap("verification_codes.Set", err)
	}

	return nil
}

func (c *VerificationCodesRepository) Delete(ctx context.Context, email string) error {
	key := c.constructKey(email)

	err := c.cache.Delete(ctx, key)
	if err != nil {
		return errs.Wrap("verification_codes.Delete", err)
	}

	return nil
}

func (c *VerificationCodesRepository) constructKey(key string) string {
	return fmt.Sprintf("%s:%s", consts.CacheVerifyCodePrefix, key)
}
