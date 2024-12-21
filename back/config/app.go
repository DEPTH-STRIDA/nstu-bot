package config

import (
	"app/log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// LoadConfig загружает конфигурацию из файла .env и переменных окружения.
//
//	configFile - ссылка на структуру, в которую будут загружены параметры
func LoadConfig(configFile interface{}) error {
	// Загрузка файла .env
	if err := godotenv.Load(); err != nil {
		return err
	}

	err := envconfig.Process("", configFile)
	if err != nil {
		return err
	}
	log.App.Info("Загруженые параметры: \n", configFile)
	return nil
}
