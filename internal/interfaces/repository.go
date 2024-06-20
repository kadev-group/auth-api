package interfaces

import (
	"auth-api/internal/models"
	"context"
)

type IRepository interface {
	Users() IUserRepository
	Sessions() ISessionRepository
	SessionsCache() ISessionsCacheRepository
	RequestSessions() IRequestSessionRepository
	ValidateCodes() IValidateCodesRepository
}

type IUserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id int64) (result *models.User, err error)
	FindByEmail(ctx context.Context, email string) (result *models.User, err error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (user *models.User, err error)
	FindByUserIDCode(ctx context.Context, userIDCode string) (user *models.User, err error)
	FindWSessionByToken(ctx context.Context, refreshToken string) (*models.UserSession, error)
}

type ISessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	FindByID(ctx context.Context, sessionID int64) (*models.Session, error)
	FindByToken(ctx context.Context, refreshToken string) (*models.Session, error)
	UpdateByID(ctx context.Context, session *models.Session) error
	UpdateByUserID(ctx context.Context, session *models.Session) error
	EndSession(ctx context.Context, sessionID int64) error
	DeleteByID(ctx context.Context, sessionID int64) error
}

type ISessionsCacheRepository interface {
	Get(ctx context.Context, userIDCode string) (int64, error)
	Set(ctx context.Context, userIDCode string, startedAt int64) error
}

type IRequestSessionRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, code string) error
}

type IValidateCodesRepository interface {
	Get(ctx context.Context, key string) (*models.VerificationCode, error)
	Set(ctx context.Context, key string, val *models.VerificationCode) error
	Delete(ctx context.Context, key string) error
}
