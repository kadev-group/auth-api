package interfaces

import (
	"auth-api/internal/models"
	"context"
)

type IService interface {
	Auth() IAuthService
	User() IUserService
	OAuth() IOAuthService
}

type IAuthService interface {
	NewPairTokens(uSession *models.UserSession) (result *models.Tokens, err error)
	NewSession(ctx context.Context, user *models.User) (result *models.Tokens, err error)
	UpdateSession(ctx context.Context, user *models.User) (result *models.Tokens, err error)
	ValidateRefreshToken(ctx context.Context, refreshToken string) (*models.UserSession, error)
	ValidateAccessToken(ctx context.Context, accessToken string) (*models.UserSession, error)
}

type IUserService interface {
	Create(ctx context.Context, userDTO *models.UserDTO) (*models.AuthResponse, error)
	Authenticate(ctx context.Context, req *models.AuthenticateReq) (*models.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*models.Tokens, error)
	Logout(ctx context.Context, refreshToken string) (err error)
	SendVerifyCode(ctx context.Context, request *models.SendVerifyCodeReq) error
	GetByUserIDCode(ctx context.Context, userIDCode string) (*models.UserDTO, error)
}

type IOAuthService interface {
	Google() IGoogleAPI
}

type IGoogleAPI interface {
	GetRedirectURL(ctx context.Context, state string) (*models.GoogleRedirectRes, error)
	HandleCallBack(ctx context.Context, code, exchangeCode string) (string, error)
}
