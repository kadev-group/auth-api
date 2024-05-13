package models

type MailsProducerMsg struct {
	SendTo           string `json:"send_to"`
	VerificationCode string `json:"verification_code"`
}
