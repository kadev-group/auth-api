package oauth

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/pkg/tools"
	"context"
	"github.com/doxanocap/pkg/errs"
	"github.com/google/uuid"
)

type GmailService struct {
	config  *models.Config
	manager interfaces.IManager
}

func InitGmailService(manager interfaces.IManager, config *models.Config) *GmailService {
	return &GmailService{
		config:  config,
		manager: manager,
	}
}

func (gs *GmailService) GetRedirectURL(ctx context.Context, state string) (response *models.GoogleRedirectRes, err error) {
	if !tools.IsUUID(state) {
		return nil, models.ErrInvalidAuthState
	}

	value, err := gs.manager.Repository().RequestSessions().Get(ctx, state)
	if err != nil {
		return
	}
	if value != "" {
		return nil, models.HttpConflict
	}
	if err = gs.manager.Repository().RequestSessions().Set(ctx, state, uuid.NewString()); err != nil {
		return
	}

	return gs.manager.Processor().APIs().GoogleAPI().NewRedirectURL(state)
}

func (gs *GmailService) HandleCallBack(ctx context.Context, state, exchangeCode string) (string, error) {
	requestID, err := gs.manager.Repository().RequestSessions().Get(ctx, state)
	if err != nil {
		return "", err
	}
	if requestID == "" {
		return "", models.ErrInvalidState
	}

	token, err := gs.manager.Processor().APIs().GoogleAPI().Exchange(ctx, exchangeCode)
	if err != nil {
		return "", errs.Wrap("gs.exchange", err)
	}

	up, err := gs.manager.Processor().APIs().GoogleAPI().GetUserProfileByToken(ctx, token.AccessToken)
	if err != nil {
		return "", errs.Wrap("gs.GetUserInfo", err)
	}

	if err = gs.manager.Repository().RequestSessions().Delete(ctx, state); err != nil {
		return "", err
	}

	user, err := gs.manager.Repository().Users().FindByEmail(ctx, up.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		tokens, err := gs.manager.Service().User().Create(ctx, &models.UserDTO{
			Email:        up.Email,
			Activated:    true,
			AuthProvider: models.GoogleAuth,
		})
		if err != nil {
			return "", err
		}

		return tokens.Tokens.AccessToken, nil
	}

	tokens, err := gs.manager.Service().Auth().UpdateSession(ctx, user)
	if err != nil {
		return "", err
	}

	return tokens.AccessToken, nil
}

func (gs *GmailService) GmailAuth(ctx context.Context, googleToken string) (response *models.GmailAuthRes, err error) {
	response = &models.GmailAuthRes{}

	up, err := gs.manager.Processor().APIs().GoogleAPI().GetUserProfileByToken(ctx, googleToken)
	if err != nil {
		return nil, err
	}

	user, err := gs.manager.Repository().Users().FindByEmail(ctx, up.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		if user.DeletedAt != nil {
			return nil, models.ErrInactiveUser
		}

		tokens, err := gs.manager.Service().Auth().NewSession(ctx, models.GoogleAuth, user)
		if err != nil {
			return nil, err
		}
		response.Data.Tokens = *tokens
		return response, nil
	}

	response.NewUser = true
	response.Data.RequestID = uuid.NewString()

	err = gs.NewRequestSession(ctx, response.Data.RequestID, up.Email)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (gs *GmailService) NewRequestSession(ctx context.Context, requestID, email string) error {
	if !tools.IsUUID(requestID) {
		return models.HttpBadRequest
	}

	return gs.manager.Repository().RequestSessions().Set(ctx, requestID, email)
}

func (gs *GmailService) ValidateRequestSession(ctx context.Context, requestID string) (string, error) {
	email, err := gs.manager.Repository().RequestSessions().Get(ctx, requestID)
	if err != nil {
		return "", err
	}
	if email == "" {
		return "", models.ErrInvalidRequestSession
	}
	return email, nil
}
