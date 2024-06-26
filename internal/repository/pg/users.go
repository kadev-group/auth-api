package pg

import (
	"auth-api/internal/models"
	"context"
	"database/sql"
	"errors"
	"github.com/doxanocap/pkg/errs"
	"github.com/jmoiron/sqlx"
)

type UsersRepository struct {
	db *sqlx.DB
}

func InitUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

// Create ...
func (repo *UsersRepository) Create(ctx context.Context, user *models.User) (err error) {
	defer errs.WrapIfErr("repo.user.Create", &err)

	err = repo.db.QueryRowxContext(ctx,
		`insert into users
		(user_idcode, email, phone_number, password, oauth_provider, created_at) 
		values ($1,$2,$3,$4,$5, $6)
		returning user_id`,
		user.IDCode, user.Email, user.PhoneNumber,
		user.Password, user.OAuthProvider, user.CreatedAt).
		Scan(&user.ID)
	return
}

func (repo *UsersRepository) FindByID(ctx context.Context, id int64) (user *models.User, err error) {
	defer errs.WrapIfErr("repo.user.FindByID", &err)

	user = &models.User{}
	err = repo.db.GetContext(ctx, user,
		`select * from users where user_id = $1`, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return
}

func (repo *UsersRepository) FindByUserIDCode(ctx context.Context, userIDCode string) (user *models.User, err error) {
	defer errs.WrapIfErr("repo.user.FindByUserIDCode", &err)

	user = &models.User{}
	err = repo.db.GetContext(ctx, user,
		`select * from users where user_idcode = $1`, userIDCode)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return
}

func (repo *UsersRepository) FindByEmail(ctx context.Context, email string) (user *models.User, err error) {
	defer errs.WrapIfErr("repo.user.FindByEmail", &err)

	user = &models.User{}
	err = repo.db.GetContext(ctx, user,
		`select * from users where email = $1 and deleted_at is null`, email)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return
}

// FindWSessionByToken ...
func (repo *UsersRepository) FindWSessionByToken(ctx context.Context, refreshToken string) (us *models.UserSession, err error) {
	defer errs.WrapIfErr("repo.user.FindByToken", &err)

	us = &models.UserSession{}
	err = repo.db.GetContext(ctx, us, `
		select 
			user_id, user_code, session_id, session_ip, 
			refresh_token, started_at, ended_at
		from users u 
		left join sessions s 
			on s.user_idref = u.user_id
		where refresh_token = $1`, refreshToken)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return
}
