package cache

import (
	"auth-api/internal/interfaces"
	"context"
	"encoding/json"
	"github.com/doxanocap/pkg/errs"
	"go.uber.org/zap"
	"time"
)

type Cache struct {
	provider interfaces.ICacheProvider
	log      *zap.Logger
}

func NewCacheProcessor(provider interfaces.ICacheProvider, log *zap.Logger) *Cache {
	return &Cache{
		provider: provider,
		log:      log,
	}
}

func (c *Cache) Set(ctx context.Context, key string, value []byte) error {
	log := c.log.With(zap.String("key", key), zap.String("value", string(value)))

	err := c.provider.Set(ctx, key, value)
	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("cache.processor.Set", err)
	}

	log.Info("set")
	return nil
}

func (c *Cache) SetJSON(ctx context.Context, key string, value interface{}) error {
	raw, err := json.Marshal(value)
	log := c.log.With(zap.String("key", key), zap.String("value", string(raw)))

	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("marshal", err)
	}

	err = c.provider.Set(ctx, key, raw)
	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("cache.processor.SetJSON", err)
	}

	log.Info("setJSON")
	return nil
}

func (c *Cache) SetJSONWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	raw, err := json.Marshal(value)
	log := c.log.With(
		zap.String("key", key),
		zap.String("value", string(raw)),
		zap.Duration("ttl", ttl))

	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("marshal", err)
	}

	err = c.provider.SetWithTTL(ctx, key, raw, ttl)
	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("cache.processor.SetJSONWithTTL", err)
	}

	log.Info("setJSONWithTTL")
	return nil
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	log := c.log.With(zap.String("key", key))

	raw, err := c.provider.Get(ctx, key)
	if err != nil {
		log.Error(err.Error())
		return nil, errs.Wrap("cache.processor.Get", err)
	}

	log.Info("get")
	return raw, nil
}

func (c *Cache) GetJSON(ctx context.Context, key string, v interface{}) error {
	log := c.log.With(zap.String("key", key))

	raw, err := c.provider.Get(ctx, key)
	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("cache.processor.GetJSON", err)
	}

	err = json.Unmarshal(raw, v)
	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("unmarshal", err)
	}

	log.Info("getJSON")
	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	log := c.log.With(zap.String("key", key))

	err := c.provider.Delete(ctx, key)
	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("cache.processor.Delete", err)
	}

	log.Info("delete")
	return nil
}

func (c *Cache) FlushAll(ctx context.Context) error {
	err := c.provider.FlushAll(ctx)
	if err != nil {
		c.log.Error(err.Error())
		return errs.Wrap("cache.processor.FlushAll", err)
	}
	c.log.Info("flushAll")
	return nil
}
