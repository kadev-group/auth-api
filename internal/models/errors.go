package models

import (
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
	ErrInvalidOAuthState  = errs.NewHttp(http.StatusBadRequest, "invalid oauth state")
	ErrInvalidRequest     = errs.NewHttp(http.StatusBadRequest, "invalid request params")
	ErrInvalidEmail       = errs.NewHttp(http.StatusBadRequest, "invalid email")
	ErrInvalidPassword    = errs.NewHttp(http.StatusBadRequest, "invalid password format")
	ErrInvalidPhoneNumber = errs.NewHttp(http.StatusBadRequest, "invalid phone number")

	ErrInvalidToken      = errs.NewHttp(http.StatusUnauthorized, "invalid token")
	ErrIncorrectPassword = errs.NewHttp(http.StatusUnauthorized, "incorrect password")

	ErrStateNotFound   = errs.NewHttp(http.StatusNotFound, "state not found")
	ErrSessionNotFound = errs.NewHttp(http.StatusNotFound, "session not found")
	ErrUserNotFound    = errs.NewHttp(http.StatusNotFound, "user not found")

	ErrUserMustAuthWGoogle = errs.NewHttp(http.StatusConflict, "user must proceed with google")
	ErrUserAlreadyExist    = errs.NewHttp(http.StatusConflict, "user already exist")
	ErrSessionExpired      = errs.NewHttp(http.StatusConflict, "session is expired")
	ErrStateCollision      = errs.NewHttp(http.StatusConflict, "such oauth state already exist")

	ErrInvalidOAuthProvider = errs.New("invalid oauth provider")
)
