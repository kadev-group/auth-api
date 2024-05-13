package repository

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/repository/cache"
	"auth-api/internal/repository/pg"
	"github.com/jmoiron/sqlx"
	"sync"
)

type Repository struct {
	db     *sqlx.DB
	config *models.Config
	cache  interfaces.ICacheProcessor

	user       interfaces.IUserRepository
	userRunner sync.Once

	sessions       interfaces.ISessionRepository
	sessionsRunner sync.Once

	sessionsCache       interfaces.ISessionsCacheRepository
	sessionsCacheRunner sync.Once

	googleAPICodes       interfaces.IGoogleAPICodesRepository
	googleAPICodesRunner sync.Once

	verificationCodeCache       interfaces.IVerificationCodesRepository
	verificationCodeCacheRunner sync.Once
}

func InitRepository(
	db *sqlx.DB,
	config *models.Config,
	cache interfaces.ICacheProcessor) *Repository {
	return &Repository{
		db:     db,
		cache:  cache,
		config: config,
	}
}

func (r *Repository) Users() interfaces.IUserRepository {
	r.userRunner.Do(func() {
		r.user = pg.InitUsersRepository(r.db)
	})
	return r.user
}

func (r *Repository) Sessions() interfaces.ISessionRepository {
	r.sessionsRunner.Do(func() {
		r.sessions = pg.InitSessionsRepository(r.db)
	})
	return r.sessions
}

func (r *Repository) SessionsCache() interfaces.ISessionsCacheRepository {
	r.sessionsCacheRunner.Do(func() {
		r.sessionsCache = cache.InitSessionCacheRepository(r.cache)
	})
	return r.sessionsCache
}

func (r *Repository) GoogleAPICodes() interfaces.IGoogleAPICodesRepository {
	r.googleAPICodesRunner.Do(func() {
		r.googleAPICodes = cache.InitOAuthCacheRepository(models.GoogleOAuth, r.cache)
	})
	return r.googleAPICodes
}

func (r *Repository) VerificationCodes() interfaces.IVerificationCodesRepository {
	r.verificationCodeCacheRunner.Do(func() {
		r.verificationCodeCache = cache.InitVerificationCodesRepository(r.cache)
	})
	return r.verificationCodeCache
}
