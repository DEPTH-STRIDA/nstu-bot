package model

import (
	"app/internal/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	WebConfig
	TelegramConfig
	DataBaseConfig
	MoyClassConfig
	SMTPConfig
}

type WebConfig struct {
	APPIP        string `envconfig:"APP_IP" default:"localhost"`    // IP адрес приложения
	APPPORT      string `envconfig:"APP_PORT" default:"8080"`       // Порт приложения
	MailPassword string `envconfig:"MAIL_PASSWORD" required:"true"` // Порт приложения
}

type TelegramConfig struct {
	Token                      string  `envconfig:"TELEGRAM_TOKEN" required:"true"`                                  // Токен бота
	Admins                     []int64 `envconfig:"TELEGRAM_ADMINS"`                                                 // Список ID админов
	RequestUpdatePause         int     `envconfig:"TELEGRAM_REQUEST_UPDATE_PAUSE__MIILISECOND" default:"3000"`       // Пауза между обработкой обновлений. Например, отправки сообщений
	RequestCallBackUpdatePause int     `envconfig:"TELEGRAM_REQUEST_CALLBACK_UPDATE_PAUSE_MIILISECOND default:"500"` // Пауза между запросами на обновление callback
	MsgBufferSize              int     `envconfig:"TELEGRAM_MESSAGE_BUFFER_SIZE" default:"100"`                      // Размер буфера для сообщений
	CallBackBufferSize         int     `envconfig:"TELEGRAM_CALLBACK_BUFFER_SIZE" default:"100"`                     // Размер буфера для callback
	TutorChatId                int64   `envconfig:"TELEGRAM_TUTOR_CHAT_ID" required:"true"`
}

type DataBaseConfig struct {
	Host     string `envconfig:"DBHOST" required:"true"` // IP адресс для подключение к БД
	Port     string `envconfig:"DBPORT" default:""`      // Port для подключение к БД
	DBName   string `envconfig:"DBNAME" required:"true"` // Имя базы данных
	UserName string `envconfig:"DBUSER" required:"true"` // Имя пользователя
	Password string `envconfig:"DBPASS" required:"true"` // Пароль пользователя
	SSLMode  string `envconfig:"DBSSLMODE" default:"disable"`
}

type MoyClassConfig struct {
	WebHouckUrl string `envconfig:"MOY_KLASS_WEBHOUK_URL" required:"true"` // Путь для обработки вебхуков мой класс
	ApiKey      string `envconfig:"API_KEY" required:"true"`               // Путь для обработки вебхуков мой класс
}

type SMTPConfig struct {
	Mail                  string `envconfig:"SMTP_MAIL" required:"true"`                  // Email адрес отправителя
	Password              string `envconfig:"SMTP_PASSWORD" required:"true"`              // Пароль от почты
	Host                  string `envconfig:"SMTP_HOST" required:"true"`                  // SMTP хост
	Port                  int    `envconfig:"SMTP_PORT" required:"true"`                  // SMTP порт с шифрованием
	PortWithoutEncryption int    `envconfig:"SMTP_PORT_WITHOUT_ENCRYPTION" default:"587"` // SMTP порт без шифрования
}

var ConfigFile *Config

func InitConfig() error {
	// Загрузка файла .env
	if err := godotenv.Load(); err != nil {
		return err
	}

	ConfigFile = &Config{}
	err := envconfig.Process("", ConfigFile)
	if err != nil {
		return err
	}
	logger.Log.Info("Загруженые параметры: \n", ConfigFile)
	return nil
}
