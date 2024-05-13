package models

import (
	"auth-api/internal/pkg/tools"
)

type SignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r SignInReq) ToUserDTO() *UserDTO {
	return &UserDTO{
		Email:    r.Email,
		Password: r.Password,
	}
}

func (r SignInReq) Validate() error {
	if !tools.IsValidEmail(r.Email) {
		return ErrInvalidEmail
	}
	if !tools.IsValidPassword(r.Password) {
		return ErrInvalidPassword
	}
	return nil
}

type SignUpReq struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

func (r SignUpReq) ToUserDTO() *UserDTO {
	return &UserDTO{
		Email:         r.Email,
		Password:      r.Password,
		PhoneNumber:   r.PhoneNumber,
		OAuthProvider: DefaultOAuth,
	}
}

func (r SignUpReq) Validate() error {
	if !tools.IsValidEmail(r.Email) {
		return ErrInvalidEmail
	}
	if !tools.IsValidPassword(r.Password) {
		return ErrInvalidPassword
	}
	if !tools.IsValidPhoneNumber(r.PhoneNumber) {
		return ErrInvalidPhoneNumber
	}
	return nil
}

type VerifyEmailReq struct {
	Email string `json:"email"`
}

func (r VerifyEmailReq) Validate() error {
	if !tools.IsValidEmail(r.Email) {
		return ErrInvalidEmail
	}
	return nil
}
