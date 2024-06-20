package models

import (
	"auth-api/internal/pkg/tools"
)

type WebSignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r WebSignInReq) Validate() error {
	if !tools.IsValidEmail(r.Email) ||
		!tools.IsValidPassword(r.Password) {
		return HttpBadRequest
	}
	return nil
}

func (r WebSignInReq) ToAuthReq() *AuthenticateReq {
	return &AuthenticateReq{
		Email:        r.Email,
		Password:     r.Password,
		AuthProvider: EmailAuth,
	}
}

type MobileSignInReq struct {
	RequestID    string `json:"request_id"`
	PhoneNumber  string `json:"phone_number"`
	ValidateCode string `json:"validate_code"`
}

func (r *MobileSignInReq) Validate() error {
	if r == nil ||
		!tools.IsUUID(r.RequestID) ||
		!tools.IsValidPhoneNumber(r.PhoneNumber) {
		return HttpBadRequest
	}

	r.PhoneNumber = tools.NormalizePhone(r.PhoneNumber)
	return nil
}

func (r *MobileSignInReq) ToAuthReq() *AuthenticateReq {
	if r == nil {
		return nil
	}
	return &AuthenticateReq{
		PhoneNumber:  r.PhoneNumber,
		RequestID:    r.RequestID,
		ValidateCode: r.ValidateCode,
		AuthProvider: PhoneAuth,
	}
}

type MobileGmailAuthReq struct {
	GoogleToken string `json:"google_access_token"`
}

type SignUpReq struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

func (r SignUpReq) ToUserDTO() *UserDTO {
	return &UserDTO{
		Email:        r.Email,
		Password:     r.Password,
		PhoneNumber:  r.PhoneNumber,
		AuthProvider: EmailAuth,
	}
}

func (r SignUpReq) Validate() error {
	if !tools.IsValidEmail(r.Email) ||
		!tools.IsValidPassword(r.Password) ||
		!tools.IsValidPhoneNumber(r.PhoneNumber) {
		return HttpBadRequest
	}
	return nil
}

type SendValidateCodeReq struct {
	Email string `json:"email"`

	RequestID   string `json:"request_id"`
	PhoneNumber string `json:"phone_number"`

	AuthProvider AuthProvider `json:"auth_provider"`
}

func (r *SendValidateCodeReq) Validate() error {
	if r == nil {
		return HttpBadRequest
	}

	if r.AuthProvider == PhoneAuth {
		if !tools.IsUUID(r.RequestID) || !tools.IsValidPhoneNumber(r.PhoneNumber) {
			return HttpBadRequest
		}
		r.PhoneNumber = tools.NormalizePhone(r.PhoneNumber)
		return nil
	} else if r.AuthProvider == EmailAuth {
		if !tools.IsValidEmail(r.Email) {
			return HttpBadRequest
		}
		return nil
	} else {
		return HttpBadRequest
	}
}
