package service

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"auth-api/internal/pkg/tools"
	"context"
	"encoding/json"
	"errors"
	"github.com/doxanocap/pkg/ctxholder"
	"github.com/doxanocap/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthService struct {
	config  *models.Config
	manager interfaces.IManager
}

func InitAuthService(manager interfaces.IManager, config *models.Config) *AuthService {
	return &AuthService{
		config:  config,
		manager: manager,
	}
}

func (s *AuthService) NewSession(
	ctx context.Context,
	provider models.AuthProvider,
	user *models.User) (result *models.Tokens, err error) {
	startedAt := time.Now().Unix()
	userSession := &models.UserSession{
		UserIDCode:    user.IDCode,
		UserEmail:     user.Email,
		UserActivated: user.Activated,
		StartedAt:     startedAt,
	}

	tokens, err := s.NewPairTokens(userSession)
	if err != nil {
		return nil, err
	}

	if err = s.manager.Repository().Sessions().Create(ctx, &models.Session{
		UserIDRef:    user.ID,
		AuthProvider: provider,
		IP:           s.getClientIP(ctx),
		RefreshToken: tokens.RefreshToken,
		StartedAt:    tools.GetPtr(startedAt),
	}); err != nil {
		return nil, err
	}

	if err = s.manager.Repository().SessionsCache().Set(ctx, user.IDCode, startedAt); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) UpdateSession(ctx context.Context, user *models.User) (result *models.Tokens, err error) {
	startedAt := time.Now().Unix()
	userSession := &models.UserSession{
		UserIDCode:    user.IDCode,
		UserEmail:     user.Email,
		UserActivated: user.Activated,
		StartedAt:     startedAt,
	}

	tokens, err := s.NewPairTokens(userSession)
	if err != nil {
		return nil, err
	}

	if err = s.manager.Repository().Sessions().UpdateByUserID(ctx, &models.Session{
		UserIDRef:    user.ID,
		RefreshToken: tokens.RefreshToken,
		IP:           s.getClientIP(ctx),
		StartedAt:    tools.GetPtr(startedAt),
	}); err != nil {
		return nil, err
	}

	if err = s.manager.Repository().SessionsCache().Set(ctx, user.IDCode, startedAt); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) NewPairTokens(userSession *models.UserSession) (result *models.Tokens, err error) {
	payload, err := json.Marshal(userSession)
	if err != nil {
		return
	}

	accessClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.AccessTokenTTL)),
		Issuer:    string(payload),
	})
	refreshClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.RefreshTokenTTL)),
		Issuer:    string(payload),
	})
	accessToken, _ := accessClaim.SignedString([]byte(s.config.Token.AccessSecret))
	refreshToken, _ := refreshClaim.SignedString([]byte(s.config.Token.RefreshSecret))

	return &models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) ValidateRefreshToken(ctx context.Context, refreshToken string) (*models.UserSession, error) {
	session, err := s.manager.Repository().Sessions().FindByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, models.ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Token.RefreshSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, models.ErrInvalidToken
		}
		return nil, errs.Wrap("parse with claims", err)
	}

	userSession := models.UserSession{UserID: session.UserIDRef}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	err = json.Unmarshal([]byte(claims.Issuer), &userSession)
	if err != nil {
		return nil, errs.Wrap("unmarshal claims", err)
	}

	if !ok || !token.Valid || userSession.StartedAt != *session.StartedAt {
		return nil, models.ErrInvalidToken
	}
	return &userSession, nil
}

func (s *AuthService) ValidateAccessToken(ctx context.Context, accessToken string) (*models.UserSession, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Token.AccessSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, models.ErrInvalidToken
		}
		return nil, errs.Wrap("parse with claims", err)
	}

	var userSession models.UserSession
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	err = json.Unmarshal([]byte(claims.Issuer), &userSession)
	if err != nil {
		return nil, errs.Wrap("unmarshal claims", err)
	}

	if !ok || !token.Valid || userSession.UserIDCode == "" {
		return nil, models.ErrInvalidToken
	}

	startedAt, err := s.manager.Repository().SessionsCache().Get(ctx, userSession.UserIDCode)
	if err != nil {
		return nil, err
	}
	if startedAt != userSession.StartedAt {
		return nil, models.ErrSessionExpired
	}

	return &userSession, nil
}

func (s *AuthService) getClientIP(ctx context.Context) (clientIP string) {
	c, ok := ctx.(*gin.Context)
	if ok {
		clientIP = ctxholder.GetClientIP(c)
	}
	return
}
