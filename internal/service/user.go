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
	smsFreePhones = map[string]string{
		"70000000000": "7456",
		"77079999999": "1111",
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

	tokens, err := us.manager.Service().Auth().NewSession(ctx, userDTO.AuthProvider, user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		UserDTO: *userDTO,
		Tokens:  tokens,
	}, nil
}

func (us *UserService) Authenticate(ctx context.Context, request *models.AuthenticateReq) (*models.AuthResponse, error) {
	var user *models.User
	var tokens *models.Tokens
	var err error

	switch request.AuthProvider {
	case models.EmailAuth:
		user, err = us.manager.Repository().Users().FindByEmail(ctx, request.Email)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, models.ErrUserNotFound
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
		if err != nil {
			return nil, models.ErrIncorrectPassword
		}

		tokens, err = us.manager.Service().Auth().UpdateSession(ctx, user)
		if err != nil {
			return nil, err
		}

		return &models.AuthResponse{
			UserDTO: user.ToUserDTO(),
			Tokens:  tokens,
		}, nil
	case models.PhoneAuth:
		email, err := us.manager.Service().OAuth().Gmail().ValidateRequestSession(ctx, request.RequestID)
		if err != nil {
			return nil, err
		}

		if err = us.ValidateCode(ctx, request.PhoneNumber, request.ValidateCode); err != nil {
			return nil, err
		}

		user, err := us.manager.Repository().Users().FindByPhoneNumber(ctx, request.PhoneNumber)
		if err != nil {
			return nil, err
		}
		if user == nil {
			user = &models.User{
				IDCode:      uuid.New().String(),
				CreatedAt:   tools.CurrTimePtr(),
				PhoneNumber: request.PhoneNumber,
				Activated:   true,
				Email:       email,
			}
			if err = us.manager.Repository().Users().Create(ctx, user); err != nil {
				return nil, err
			}
		}

		tokens, err = us.manager.Service().Auth().NewSession(ctx, models.PhoneAuth, user)
		if err != nil {
			return nil, err
		}

		return &models.AuthResponse{
			UserDTO: user.ToUserDTO(),
			Tokens:  tokens,
		}, nil
	}
	return nil, models.ErrInvalidAuthProvider
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
	session, err := us.manager.Repository().
		Sessions().FindByToken(ctx, refreshToken)
	if err != nil {
		return err
	}
	if session == nil {
		return models.ErrInvalidToken
	}
	if err = us.manager.Repository().
		Sessions().EndSession(ctx, session.ID); err != nil {
		return err
	}
	return nil
}

func (us *UserService) SendValidateCode(ctx context.Context, request *models.SendValidateCodeReq) error {
	code := tools.NewVerificationCode()

	switch request.AuthProvider {
	case models.EmailAuth:
		user, err := us.manager.Repository().Users().FindByEmail(ctx, request.Email)
		if err != nil {
			return err
		}
		if user == nil {
			return models.ErrUserNotFound
		}

		if err = us.SetValidateCode(ctx, request.Email, code); err != nil {
			return err
		}

		if err = us.manager.Processor().Queue().Producers().Mails().Send(ctx, &models.MailsProducerMsg{
			SendTo:           request.Email,
			VerificationCode: code,
		}); err != nil {
			return err
		}
	case models.PhoneAuth:
		email, err := us.manager.Service().OAuth().Gmail().ValidateRequestSession(ctx, request.RequestID)
		if err != nil {
			return err
		}

		user, err := us.manager.Repository().Users().FindByPhoneNumber(ctx, request.PhoneNumber)
		if err != nil {
			return err
		}
		if user != nil {
			if user.DeletedAt != nil {
				return models.ErrInactiveUser
			}
			if user.Email != email {
				return models.ErrGmailAlreadyRegistered(user.Email)
			}
		}

		v, ok := smsFreePhones[request.PhoneNumber]
		if ok {
			if err = us.SetValidateCode(ctx, request.PhoneNumber, v); err != nil {
				return err
			}
			return nil
		}

		if err = us.SetValidateCode(ctx, request.PhoneNumber, code); err != nil {
			return err
		}

		if _, ok := smsFreePhones[request.PhoneNumber]; ok {
			return nil
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

func (us *UserService) SetValidateCode(ctx context.Context, phoneNumber, code string) error {
	validateCodes, err := us.manager.Repository().ValidateCodes().Get(ctx, phoneNumber)
	if err != nil {
		return err
	}

	if validateCodes.IsLimitReached() || validateCodes.IsFrequent() {
		return models.ErrValidateCodesLimit
	}

	validateCodes.Insert(code)
	if err = us.manager.Repository().ValidateCodes().
		Set(ctx, phoneNumber, validateCodes); err != nil {
		return err
	}
	return nil
}

func (us *UserService) ValidateCode(ctx context.Context, phoneNumber, code string) error {
	validateCodes, err := us.manager.Repository().ValidateCodes().Get(ctx, phoneNumber)
	if err != nil {
		return err
	}

	if !validateCodes.Find(code) {
		return models.ErrIncorrectValidateCode
	}
	return nil
}
