package models

type (
	OAuthProvider string
)

var (
	DefaultOAuth OAuthProvider = "default"
	GoogleOAuth  OAuthProvider = "google"
)

func (p OAuthProvider) IsValid() bool {
	return p == GoogleOAuth
}
