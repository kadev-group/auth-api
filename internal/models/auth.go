package models

import "time"

type AuthenticateReq struct {
	Email    string `json:"email"`
	Password string `json:"-"`

	RequestID   string `json:"request_id"`
	PhoneNumber string `json:"phone_number"`
	VerifyCode  string `json:"verify_code"`

	AuthProvider AuthProvider `json:"auth_provider"`
}

type Tokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	IssuedAt     time.Time `json:"-"`
}
