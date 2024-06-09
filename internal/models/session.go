package models

import "time"

// TODO: delete json

// Session ...
type Session struct {
	ID           int64  `json:"session_id" db:"session_id"`
	IP           string `json:"session_ip" db:"session_ip"`
	UserIDRef    int64  `json:"user_idref" db:"user_idref"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
	StartedAt    *int64 `json:"started_at" db:"started_at"`
	EndedAt      *int64 `json:"ended_at" db:"ended_at"`
}

type VerificationCode struct {
	codes        []string
	lastSentTime time.Time
}

func (p VerificationCode) IsLimitReached() bool {
	return len(p.codes) >= 2
}

func (p VerificationCode) IsFrequent() bool {
	diff := time.Now().Sub(p.lastSentTime)
	return diff < 1*time.Minute
}

func (p VerificationCode) Insert(code string) {
	p.codes = append(p.codes, code)
}
