package models

type AuthResponse struct {
	UserDTO `json:"user"`
	Tokens  *Tokens `json:"tokens"`
}

type GoogleRedirectRes struct {
	RedirectURL string `json:"redirect_url"`
}

type GmailAuthRes struct {
	NewUser bool `json:"new_user"`
	Data    struct {
		Tokens
		RequestID string `json:"request_id,omitempty"`
	} `json:"data"`
}
