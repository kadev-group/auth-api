package models

import "time"

type AuthenticateReq struct {
	Email    string `json:"email"`
	Password string `json:"-"`

	RequestID    string `json:"request_id"`
	PhoneNumber  string `json:"phone_number"`
	ValidateCode string `json:"validate_code"`

	AuthProvider AuthProvider `json:"auth_provider"`
}

type Tokens struct {
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	IssuedAt     time.Time `json:"-"`
}

type UserProfile struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}
