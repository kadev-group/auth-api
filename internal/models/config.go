package models

type Config struct {
	ENV string `env:"APP_ENV"`

	ServerPORT string `env:"SERVER_PORT"`

	PSQL
	Token
	Redis
	OAuth
	RabbitMQ
}

func (c Config) Env() string {
	return c.ENV
}

type Token struct {
	RefreshSecret string `env:"REFRESH_TOKEN_SECRET"`
	AccessSecret  string `env:"ACCESS_TOKEN_SECRET"`
}

type PSQL struct {
	PqHOST     string `env:"PSQL_HOST"`
	PqPORT     string `env:"PSQL_PORT"`
	PqUSER     string `env:"PSQL_USER"`
	PqPASSWORD string `env:"PSQL_PASSWORD"`
	PqDATABASE string `env:"PSQL_DATABASE"`
	PqSSL      string `env:"PSQL_SSL"`
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
