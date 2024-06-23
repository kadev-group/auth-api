package models

import (
	"fmt"
	"net/http"

	"github.com/doxanocap/pkg/errs"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

// default http error responses
var (
	HttpBadRequest          = errs.NewHttp(http.StatusBadRequest, "bad request")
	HttpNotFound            = errs.NewHttp(http.StatusNotFound, "not found")
	HttpInternalServerError = errs.NewHttp(http.StatusInternalServerError, "internal server error")
	HttpConflict            = errs.NewHttp(http.StatusConflict, "conflict")
	HttpUnauthorized        = errs.NewHttp(http.StatusUnauthorized, "unauthorized")
)

// custom errors for special cases
var (
	ErrInvalidAuthState    = errs.NewHttp(http.StatusBadRequest, "invalid auth state")
	ErrInvalidAuthProvider = errs.NewHttp(http.StatusBadRequest, "invalid auth provider")

	ErrInvalidToken      = errs.NewHttp(http.StatusUnauthorized, "invalid token")
	ErrSessionExpired    = errs.NewHttp(http.StatusUnauthorized, "session is expired")
	ErrIncorrectPassword = errs.NewHttp(http.StatusUnauthorized, "incorrect password")

	ErrStateNotFound   = errs.NewHttp(http.StatusNotFound, "state not found")
	ErrSessionNotFound = errs.NewHttp(http.StatusNotFound, "session not found")
	ErrUserNotFound    = errs.NewHttp(http.StatusNotFound, "user not found")

	ErrValidateCodesLimit    = errs.NewHttp(http.StatusConflict, "verification codes limit")
	ErrIncorrectValidateCode = errs.NewHttp(http.StatusConflict, "incorrect validate code")

	ErrUserMustAuthWGoogle   = errs.NewHttp(http.StatusConflict, "user must proceed with google")
	ErrUserAlreadyExist      = errs.NewHttp(http.StatusConflict, "user already exist")
	ErrInvalidState          = errs.NewHttp(http.StatusConflict, "invalid state")
	ErrInvalidRequestSession = errs.NewHttp(http.StatusConflict, "invalid request session")

	ErrInactiveUser = errs.NewHttp(http.StatusUnprocessableEntity, "user is inactive")
)

func ErrGmailAlreadyRegistered(gmail string) error {
	idx := len(gmail) - 1
	for i := idx; i >= 0; i-- {
		if gmail[i] == '@' {
			idx = i - 1
			break
		}
	}
	if idx == 0 {
		return errs.NewHttp(http.StatusInternalServerError, "invalid gmail")
	}

	hiddenGmail := []rune(gmail)

	if idx >= 8 {
		for i := idx; i > 5; i-- {
			hiddenGmail[i] = '*'
		}
	} else {
		for i := idx; i > idx/2; i-- {
			hiddenGmail[i] = '*'
		}
	}

	msg := fmt.Sprintf("Ваш номер зарегистрирован по адресу %s", string(hiddenGmail))
	return errs.NewHttp(http.StatusConflict, msg)
}
