package repository

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/repository/cache"
	"auth-api/internal/repository/pg"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	user            interfaces.IUserRepository
	sessions        interfaces.ISessionRepository
	sessionsCache   interfaces.ISessionsCacheRepository
	requestSessions interfaces.IRequestSessionRepository
	validateCodes   interfaces.IValidateCodesRepository
}

func InitRepository(
	db *sqlx.DB,
	config *models.Config,
	cacheProcessor interfaces.ICacheProcessor) *Repository {
	return &Repository{
		user:            pg.InitUsersRepository(db),
		sessions:        pg.InitSessionsRepository(db),
		sessionsCache:   cache.InitSessionCacheRepository(cacheProcessor),
		validateCodes:   cache.InitValidateCodesRepository(cacheProcessor),
		requestSessions: cache.InitRequestSessionRepository(cacheProcessor),
	}
}

func (r *Repository) Users() interfaces.IUserRepository {
	return r.user
}

func (r *Repository) Sessions() interfaces.ISessionRepository {
	return r.sessions
}

func (r *Repository) SessionsCache() interfaces.ISessionsCacheRepository {
	return r.sessionsCache
}

func (r *Repository) RequestSessions() interfaces.IRequestSessionRepository {
	return r.requestSessions
}

func (r *Repository) ValidateCodes() interfaces.IValidateCodesRepository {
	return r.validateCodes
}
