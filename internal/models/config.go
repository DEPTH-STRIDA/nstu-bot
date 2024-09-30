package models

import (
	"nstu/internal/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Config *ConfigStruct

type ConfigStruct struct {
	WebAppConfig
	TelegramConfig
	DataBaseConfig
}

type WebAppConfig struct {
	AppIP   string `envconfig:"APP_API" default:"localhost"`
	AppPort string `envconfig:"APP_PORT" default:"8000"`
}

type TelegramConfig struct {
	NSTUToken          string `envconfig:"NSTU_TOKEN" required:"true"`
	MessagePauseSec    uint   `envconfig:"MSG_PAUSE" default:"3"`
	CallBackPauseSec   uint   `envconfig:"CALLBACK_PAUSE" default:"1"`
	MessageBufferSize  uint   `envconfig:"MSG_BUFFER" default:"100"`
	CallBackBufferSize uint   `envconfig:"CALLBACK_BUFFER" default:"100"`
}

type DataBaseConfig struct {
	UserName     string `envconfig:"DBUSER" required:"true"`
	Password     string `envconfig:"DBPASS" required:"true"`
	Host         string `envconfig:"DBHOST" required:"true"`
	Port         string `envconfig:"DBPORT" required:"true"`
	DataBaseName string `envconfig:"DBNAME" required:"true"`
}

func InitConfig() error {
	// Загрузка файла .env
	if err := godotenv.Load(); err != nil {
		return err
	}

	Config = &ConfigStruct{}
	err := envconfig.Process("", Config)
	if err != nil {
		return err
	}
	logger.Log.Info("Загруженые параметры: \n", Config)
	return nil
}
