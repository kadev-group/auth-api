package interfaces

import (
	"auth-api/internal/models"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/oauth2"
	"time"
)

type IProcessor interface {
	Cache() ICacheProcessor
	Queue() IQueueProcessor
	APIs() IAPIsProcessor
	SMS() ISMSProcessor
}

// Cache

type ICacheProcessor interface {
	Set(ctx context.Context, key string, value []byte) error
	SetJSON(ctx context.Context, key string, value interface{}) error
	SetWithTTL(ctx context.Context, key string, value []byte, ttl time.Duration) error
	SetJSONWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	GetJSON(ctx context.Context, key string, v interface{}) error
	Delete(ctx context.Context, key string) error
	FlushAll(ctx context.Context) error
}

type ICacheProvider interface {
	Set(ctx context.Context, key string, value []byte) error
	SetWithTTL(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	FlushAll(ctx context.Context) error
	Close() error
}

// Queue

type IQueueProcessor interface {
	Producers() IQueueProducersProcessor
}

type IQueueProducersProcessor interface {
	Mails() IQueueProducerProcessor
}

type IQueueProducerProcessor interface {
	Send(ctx context.Context, message *models.MailsProducerMsg) error
}

type IQueueProducerProvider interface {
	Send(ctx context.Context, qName string, message []byte, args ...amqp.Table) (err error)
}

// SMS

type ISMSProcessor interface {
	Send(ctx context.Context, phone string, message string) error
}

type ISMSProvider interface {
	Send(ctx context.Context, phone string, message string) error
}

type IAPIsProcessor interface {
	GoogleAPI() IGoogleAPIProcessor
}

type IGoogleAPIProcessor interface {
	GetUserProfileByToken(ctx context.Context, accessToken string) (*models.UserProfile, error)
	NewRedirectURL(state string) (res *models.GoogleRedirectRes, err error)
	Exchange(ctx context.Context, exchangeCode string) (token *oauth2.Token, err error)
}
