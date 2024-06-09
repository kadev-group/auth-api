package models

import "time"

// User ...
type User struct {
	ID           int64        `db:"user_id"`
	IDCode       string       `db:"user_idcode"`
	Email        string       `db:"email"`
	PhoneNumber  string       `db:"phone_number"`
	Activated    bool         `db:"activated"`
	Password     string       `db:"password"`
	AuthProvider AuthProvider `db:"auth_provider"`
	CreatedAt    *time.Time   `db:"created_at"`
	UpdatedAt    *time.Time   `db:"updated_at"`
	DeletedAt    *time.Time   `db:"deleted_at"`
}

// ToUserDTO ...
func (u *User) ToUserDTO() UserDTO {
	return UserDTO{
		IDCode:      u.IDCode,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Activated:   u.Activated,
		Password:    u.Password,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		DeletedAt:   u.DeletedAt,
	}
}

// UserDTO ...
type UserDTO struct {
	IDCode       string       `json:"id"`
	Email        string       `json:"email"`
	PhoneNumber  string       `json:"phone_number"`
	Activated    bool         `json:"activated"`
	Password     string       `json:"-"`
	AuthProvider AuthProvider `json:"auth_provider"`
	CreatedAt    *time.Time   `json:"created_at"`
	UpdatedAt    *time.Time   `json:"updated_at"`
	DeletedAt    *time.Time   `json:"deleted_at"`
}

// ToUser ...
func (u *UserDTO) ToUser() *User {
	return &User{
		IDCode:       u.IDCode,
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		Activated:    u.Activated,
		Password:     u.Password,
		AuthProvider: u.AuthProvider,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		DeletedAt:    u.DeletedAt,
	}
}

type UserSession struct {
	UserID        int64  `json:"-" db:"user_id"`
	UserIDCode    string `json:"user_id" db:"user_idcode"`
	UserEmail     string `json:"email" db:"email"`
	UserActivated bool   `json:"activated" db:"activated"`
	StartedAt     int64  `json:"session_started_at" db:"started_at"`
}

// ToUser ...
func (us *UserSession) ToUser() *User {
	return &User{
		ID:        us.UserID,
		IDCode:    us.UserIDCode,
		Email:     us.UserEmail,
		Activated: us.UserActivated,
	}
}
