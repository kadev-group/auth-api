package service

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"auth-api/internal/pkg/tools"
	"context"
	"fmt"
	"github.com/doxanocap/pkg/errs"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	smsFreePhones = map[string]int{
		"70000000000": 7456,
		"77079999999": 1111,
	}
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

	if userDTO.AuthProvider == models.EmailAuth {
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

func (us *UserService) Authenticate(ctx context.Context, req *models.AuthenticateReq) (*models.AuthResponse, error) {
	var user *models.User
	var err error

	switch req.AuthProvider {
	case models.EmailAuth:
		user, err = us.manager.Repository().Users().FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, models.ErrUserNotFound
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			return nil, models.ErrIncorrectPassword
		}
	case models.PhoneAuth:
		
	case models.GoogleAuth:
		return nil, models.ErrUserMustAuthWGoogle
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

func (us *UserService) SendVerifyCode(ctx context.Context, request *models.SendVerifyCodeReq) error {
	code := tools.NewVerificationCode()

	switch request.AuthProvider {
	case models.EmailAuth:
		user, err := us.manager.
			Repository().Users().FindByEmail(ctx, request.Email)
		if err != nil {
			return err
		}
		if user == nil {
			return models.ErrUserNotFound
		}

		verifyCodes, err := us.manager.Repository().VerificationCodes().Get(ctx, request.Email)
		if err != nil {
			return err
		}

		if verifyCodes.IsLimitReached() || verifyCodes.IsFrequent() {
			return models.ErrVerifyCodesLimit
		}

		verifyCodes.Insert(code)
		if err = us.manager.Repository().VerificationCodes().
			Set(ctx, request.Email, verifyCodes); err != nil {
			return err
		}

		if err = us.manager.
			Processor().Queue().Producers().Mails().
			Send(ctx, &models.MailsProducerMsg{
				SendTo:           request.Email,
				VerificationCode: code,
			}); err != nil {
			return err
		}
	case models.PhoneAuth:
		user, err := us.manager.
			Repository().Users().FindByPhoneNumber(ctx, request.PhoneNumber)
		if err != nil {
			return err
		}
		if user == nil {
			return models.ErrUserNotFound
		}

		verifyCodes, err := us.manager.Repository().VerificationCodes().Get(ctx, request.PhoneNumber)
		if err != nil {
			return err
		}

		if verifyCodes.IsLimitReached() || verifyCodes.IsFrequent() {
			return models.ErrVerifyCodesLimit
		}

		verifyCodes.Insert(code)
		if err = us.manager.Repository().VerificationCodes().
			Set(ctx, request.PhoneNumber, verifyCodes); err != nil {
			return err
		}

		smsText := fmt.Sprintf("Ваш код подтверждения для %s - %s", consts.AppName, code)

		err = us.manager.Processor().SMS().Send(ctx, request.PhoneNumber, smsText)
		if err != nil {
			return err
		}
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
