package tools

import (
	"auth-api/internal/models/consts"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var (
	phoneNumberRegexpFn = regexp.MustCompile(consts.PhoneNumberRegexp)
	emailRegexpFn       = regexp.MustCompile(consts.EmailRegexp)
)

func CurrTimePtr() *time.Time {
	t := time.Now()
	return &t
}

func GetPtr[T any](v T) *T {
	return &v
}

func IsValidVerificationCode(code string) bool {
	if len(code) > 6 {
		return false
	}
	for i := range code {
		if code[i] < 48 || code[i] > 57 {
			return false
		}
	}
	return true
}

func NewVerificationCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d", r.Intn(900000)+100000)
}

func IsValidDateFormat(date string) bool {
	if len(date) != len(consts.DateFormat) {
		return false
	}
	split := strings.Split(date, "-")
	if len(split) != 3 {
		return false
	}

	return true
}

func IsValidPassword(password string) bool {
	return len(password) > 8 && len(password) < 48
}

func IsUUID(str string) bool {
	_, err := uuid.Parse(str)
	return err == nil
}

func IsValidEmail(e string) bool {
	return emailRegexpFn.MatchString(e)
}

func IsValidPhoneNumber(e string) bool {
	return phoneNumberRegexpFn.MatchString(e)
}

func NormalizePhone(p string) string {
	l := len(p)
	if l > 1 {
		if p[0] == '+' {
			p = p[1:]
		} else {
			if l == 10 && p[0] == '7' {
				p = "7" + p
			} else if l == 11 && strings.HasPrefix(p, "87") {
				p = "7" + p[1:]
			}
		}
	}
	return p
}
