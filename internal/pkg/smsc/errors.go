package smsc

import (
	"fmt"
)

var (
	ErrInvalidParameter   = &Error{Code: 1}
	ErrInvalidCredentials = &Error{Code: 2}
	ErrNoMoney            = &Error{Code: 3}
	ErrIPBlocked          = &Error{Code: 4}
	ErrInvalidDateFormat  = &Error{Code: 5}
	ErrForbidden          = &Error{Code: 6}
	ErrInvalidPhoneNumber = &Error{Code: 7}
	ErrNotDelivered       = &Error{Code: 8}
	ErrConcurrentRequests = &Error{Code: 9}

	possibleErrors = []*Error{
		ErrInvalidParameter,
		ErrInvalidCredentials,
		ErrNoMoney,
		ErrIPBlocked,
		ErrInvalidDateFormat,
		ErrForbidden,
		ErrInvalidPhoneNumber,
		ErrNotDelivered,
		ErrConcurrentRequests,
	}
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("fail: error_code(%d), error(%s)", e.Code, e.Message)
}

func Parse(code int, message string) *Error {
	for _, err := range possibleErrors {
		if err.Code == code {
			err.Message = message
			return err
		}
	}

	return nil
}
