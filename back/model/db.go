package model

type DataBaseConfig struct {
	Host     string `envconfig:"DBHOST" required:"true"` // IP адресс для подключение к БД
	Port     string `envconfig:"DBPORT" default:""`      // Port для подключение к БД
	DBName   string `envconfig:"DBNAME" required:"true"` // Имя базы данных
	UserName string `envconfig:"DBUSER" required:"true"` // Имя пользователя
	Password string `envconfig:"DBPASS" required:"true"` // Пароль пользователя
	SSLMode  string `envconfig:"DBSSLMODE" default:"disable"`
}
