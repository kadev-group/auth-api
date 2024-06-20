package models

type Config struct {
	ENV        string `env:"APP_ENV"`
	ServerPORT string `env:"SERVER_PORT"`
	PsqlDsn    string `env:"PSQL_DSN"`

	Token    Token
	SMSc     SMSc
	Redis    Redis
	OAuth    OAuth
	RabbitMQ RabbitMQ
}

func (c Config) Env() string {
	return c.ENV
}

type Token struct {
	RefreshSecret string `env:"REFRESH_TOKEN_SECRET"`
	AccessSecret  string `env:"ACCESS_TOKEN_SECRET"`
}

type Redis struct {
	Host     string `env:"REDIS_HOST"`
	Password string `env:"REDIS_PASSWORD"`
	Database int    `env:"REDIS_DATABASE"`
}

type OAuth struct {
	GoogleAPI
}

type GoogleAPI struct {
	ClientID          string `env:"OAUTH_CLIENT_ID"`
	ClientSecret      string `env:"OAUTH_CLIENT_SECRET"`
	ServerCallBackURI string `env:"OAUTH_SERVER_CALL_BACK_URI"`
	ClientCallBackURI string `env:"OAUTH_CLIENT_CALL_BACK_URI"`
}

type RabbitMQ struct {
	ServerURL  string `env:"RABBITMQ_SERVER_URL"`
	MailsQueue string `env:"RABBITMQ_MAILS_QUEUE"`
}

type SMSc struct {
	Password string `env:"SMSC_PASSWORD"`
	Login    string `env:"SMSC_LOGIN"`
}
