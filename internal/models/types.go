package models

type (
	AuthProvider string
)

var (
	EmailAuth  AuthProvider = "email"
	PhoneAuth  AuthProvider = "phone"
	GoogleAuth AuthProvider = "google"
)

func (p AuthProvider) IsValid() bool {
	return p == GoogleAuth || p == EmailAuth || p == PhoneAuth
}
