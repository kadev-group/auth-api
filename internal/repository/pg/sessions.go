package pg

import (
	"auth-api/internal/models"
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/doxanocap/pkg/errs"
	"github.com/jmoiron/sqlx"
	"time"
)

type SessionsRepository struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
}

func InitSessionsRepository(db *sqlx.DB) *SessionsRepository {
	return &SessionsRepository{
		db: db,
	}
}

// Create ...
func (repo *SessionsRepository) Create(ctx context.Context, session *models.Session) error {
	err := repo.db.QueryRowxContext(ctx, `
		insert into sessions
		(user_idref, refresh_token, session_ip, auth_provider, started_at) 
		values ($1,$2,$3,$4, $5)
		returning session_id`,
		session.UserIDRef, session.RefreshToken, session.IP,
		session.AuthProvider, session.StartedAt).
		Scan(&session.ID)
	if err != nil {
		return errs.Wrap("repository.session.Create", err)
	}
	return nil
}

// FindByID ...
func (repo *SessionsRepository) FindByID(ctx context.Context, sessionID int64) (*models.Session, error) {
	session := &models.Session{}
	err := repo.db.GetContext(ctx, session,
		`select * from sessions where session_id = $1`, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.Wrap("repository.session.FindByID", err)
	}
	return session, nil
}

// FindByToken ...
func (repo *SessionsRepository) FindByToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	session := &models.Session{}
	err := repo.db.GetContext(ctx, session,
		`select * from sessions where refresh_token = $1`, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.Wrap("repository.session.FindByToken", err)
	}
	return session, nil
}

// UpdateByID ...
func (repo *SessionsRepository) UpdateByID(ctx context.Context, session *models.Session) error {
	_, err := repo.db.ExecContext(ctx, `
		update sessions
		set 
			session_ip = $1,
			refresh_token = $2,
			started_at = $3,
			ended_at = $4
		where session_id = $5`,
		session.IP, session.RefreshToken, session.StartedAt,
		session.EndedAt, session.ID)
	if err != nil {
		return errs.Wrap("repository.session.UpdateByID", err)
	}
	return nil
}

// UpdateByUserID ...
func (repo *SessionsRepository) UpdateByUserID(ctx context.Context, session *models.Session) error {
	_, err := repo.db.ExecContext(ctx, `
		update sessions
		set 
			session_ip = $1,
			refresh_token = $2,
			started_at = $3,
			ended_at = $4
		where user_idref = $5`,
		session.IP, session.RefreshToken,
		session.StartedAt, session.EndedAt,
		session.UserIDRef)
	if err != nil {
		return errs.Wrap("repository.session.UpdateByUserID", err)
	}
	return nil
}

// EndSession ...
func (repo *SessionsRepository) EndSession(ctx context.Context, sessionID int64) error {
	_, err := repo.db.ExecContext(ctx, `
		update sessions
		set 
			session_ip = '',
			refresh_token = '',
			ended_at = $1
		where session_id = $2`,
		time.Now().Unix(), sessionID)
	if err != nil {
		return errs.Wrap("repository.session.EndSession", err)
	}
	return nil
}

// DeleteByID ...
func (repo *SessionsRepository) DeleteByID(ctx context.Context, sessionID int64) error {
	_, err := repo.db.ExecContext(ctx, `
		delete from sessions where session_id = $1`, sessionID)
	if err != nil {
		return errs.Wrap("repository.session.DeleteByID", err)
	}
	return nil
}
