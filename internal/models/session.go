package models

import (
	"time"
)

// TODO: delete json

// Session ...
type Session struct {
	ID           int64        `json:"session_id" db:"session_id"`
	IP           string       `json:"session_ip" db:"session_ip"`
	UserIDRef    int64        `json:"user_idref" db:"user_idref"`
	RefreshToken string       `json:"refresh_token" db:"refresh_token"`
	AuthProvider AuthProvider `json:"auth_provider" db:"auth_provider"`
	StartedAt    *int64       `json:"started_at" db:"started_at"`
	EndedAt      *int64       `json:"ended_at" db:"ended_at"`
}

type VerificationCode struct {
	Codes        []string  `json:"codes"`
	LastSentTime time.Time `json:"last_sent_time"`
}

func (p *VerificationCode) IsLimitReached() bool {
	return len(p.Codes) >= 2
}

func (p *VerificationCode) IsFrequent() bool {
	diff := time.Now().Sub(p.LastSentTime)
	return diff < 1*time.Minute
}

func (p *VerificationCode) Insert(code string) {
	p.Codes = append(p.Codes, code)
}

func (p *VerificationCode) Find(code string) bool {
	for i := range p.Codes {
		if p.Codes[i] == code {
			return true
		}
	}
	return false
}
