package model

type SmtpConfig struct {
	Server      string `envconfig:"SMTP_SERVER" required:"true"`
	SSLPort     int    `envconfig:"SMTP_SSL_PORT" required:"true"`
	Password    string `envconfig:"SMTP_PASSWORD" required:"true"`
	RepeatPause int    `envconfig:"SMTP_REPEAT_PAUSE" default:"1000"`
	MailName    string `envconfig:"SMTP_MAIL_NAME" required:"true"`
}
