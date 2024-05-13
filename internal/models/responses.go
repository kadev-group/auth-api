package models

type AuthResponse struct {
	UserDTO `json:"user"`
	Tokens  *Tokens `json:"tokens"`
}

type GoogleRedirectRes struct {
	RedirectURL string `json:"redirect_url"`
}
