package models

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
