package config

import "app/model"

type Config struct {
	model.WebConfig
	model.DataBaseConfig
	model.SmtpConfig
	model.AuthConfig
}

var File *Config = &Config{}

func Init() error {
	return LoadConfig(File)
}
