package smtp

import (
	"app/config"
	"app/log"
	"time"
)

var App *SMTPClient

func Init() error {
	conf := config.File.SmtpConfig
	repetPause := time.Millisecond * time.Duration(conf.RepeatPause)

	app, err := NewSMTPClient(conf.Server, conf.SSLPort, conf.MailName, conf.Password, repetPause)
	if err != nil {
		log.App.Error("failed to init smtp: ", err)
		return err
	}
	App = app
	return nil
}
