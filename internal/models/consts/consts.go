package consts

import (
	"time"
)

// App constants
const (
	AppName = "Hitba"

	EnvProd  = "prod"
	EnvDev   = "dev"
	EnvStage = "stage"

	RefreshTokenKey = "refresh_token"

	RefreshTokenTTL      = 30 * 24 * time.Hour
	AccessTokenTTL       = 30 * time.Minute
	OAuthCodeTTl         = time.Minute
	VerificationCodesTTL = 5 * time.Minute

	AuthHashCost = 10

	GoogleAwaitTime        = 5 * time.Minute
	GoogleScopeEmail       = "https://www.googleapis.com/auth/userinfo.email"
	GoogleScopeUserProfile = "https://www.googleapis.com/auth/userinfo.profile"

	CacheSessionsPrefix    = "ses"
	CacheGoogleOAuthPrefix = "oauth2"
	CacheVerifyCodePrefix  = "verify"

	DateFormat        = "2006-01-02"
	EmailRegexp       = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	PhoneNumberRegexp = `^((8|\+7)[\- ]?)?(\(?\d{3}\)?[\- ]?)?[\d\- ]{7,10}$`
)
