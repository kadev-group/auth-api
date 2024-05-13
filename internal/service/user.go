package service

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"auth-api/internal/pkg/tools"
	"context"
	"github.com/doxanocap/pkg/errs"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	manager interfaces.IManager
}

func InitUserService(manager interfaces.IManager) *UserService {
	return &UserService{
		manager: manager,
	}
}

func (us *UserService) Create(ctx context.Context, userDTO *models.UserDTO) (result *models.AuthResponse, err error) {
	found, err := us.manager.Repository().Users().FindByEmail(ctx, userDTO.Email)
	if err != nil {
		return
	}
	if found != nil {
		return nil, models.ErrUserAlreadyExist
	}

	if userDTO.OAuthProvider == models.DefaultOAuth {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), consts.AuthHashCost)
		if err != nil {
			return nil, errs.Wrap("generate hash password", err)
		}
		userDTO.Password = string(hashedPassword)
	}

	userDTO.IDCode = uuid.New().String()
	userDTO.CreatedAt = tools.CurrTimePtr()

	user := userDTO.ToUser()
	if err = us.manager.Repository().Users().Create(ctx, user); err != nil {
		return
	}

	tokens, err := us.manager.Service().Auth().NewSession(ctx, user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		UserDTO: *userDTO,
		Tokens:  tokens,
	}, nil
}

func (us *UserService) Authenticate(ctx context.Context, userDTO *models.UserDTO) (*models.AuthResponse, error) {
	user, err := us.manager.Repository().Users().FindByEmail(ctx, userDTO.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, models.ErrUserNotFound
	}

	if user.OAuthProvider == models.GoogleOAuth {
		return nil, models.ErrUserMustAuthWGoogle
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDTO.Password))
	if err != nil {
		return nil, models.ErrIncorrectPassword
	}

	tokens, err := us.manager.Service().Auth().UpdateSession(ctx, user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		UserDTO: user.ToUserDTO(),
		Tokens:  tokens,
	}, nil
}

func (us *UserService) Refresh(ctx context.Context, refreshToken string) (*models.Tokens, error) {
	userSession, err := us.manager.Service().Auth().ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	tokens, err := us.manager.Service().Auth().UpdateSession(ctx, userSession.ToUser())
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (us *UserService) Logout(ctx context.Context, refreshToken string) error {
	session, err := us.manager.Repository().Sessions().FindByToken(ctx, refreshToken)
	if err != nil {
		return err
	}
	if session == nil {
		return models.ErrInvalidToken
	}
	if err = us.manager.Repository().Sessions().EndSession(ctx, session.ID); err != nil {
		return err
	}
	return nil
}

func (us *UserService) SendVerifyCode(ctx context.Context, email string) error {
	code := tools.NewVerificationCode()

	if err := us.manager.
		Repository().
		VerificationCodes().
		Set(ctx, email, code); err != nil {
		return err
	}

	if err := us.manager.
		Processor().
		Queue().
		Producers().
		Mails().
		Send(ctx, &models.MailsProducerMsg{
			SendTo:           email,
			VerificationCode: code,
		}); err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetByUserIDCode(ctx context.Context, userIDCode string) (*models.UserDTO, error) {
	user, err := us.manager.Repository().Users().FindByUserIDCode(ctx, userIDCode)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, models.ErrUserNotFound
	}
	userDTO := user.ToUserDTO()
	return &userDTO, nil
}
