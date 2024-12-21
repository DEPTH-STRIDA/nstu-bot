package db

import (
	"app/config"
	"app/log"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

type DataBase struct {
	*gorm.DB
}

// NewDatabase создает новое подключение к базе данных
func NewDatabase() (*DataBase, error) {
	conf := config.File.DataBaseConfig
	dsn := ""
	if config.File.DataBaseConfig.Port == "" {
		dsn = fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s", conf.Host, conf.UserName, conf.DBName, conf.Password, conf.SSLMode)
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", conf.Host, conf.Port, conf.UserName, conf.DBName, conf.Password, conf.SSLMode)
	}

	// Устанавливаем логгер на Silent уровень
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.App.Info("Подключение к БД успешно установлено.")
	return &DataBase{DB: db}, nil
}
func (db *DataBase) GetRecordsByColumn(model interface{}, columnName string, limit, offset int, result interface{}) error {
	stmt := &gorm.Statement{DB: db.DB}
	_ = stmt.Parse(model)
	tableName := stmt.Schema.Table

	query := fmt.Sprintf(`SELECT * FROM "%s" ORDER BY "%s" LIMIT ? OFFSET ?`, tableName, columnName)
	dbResult := db.Raw(query, limit, offset).Scan(result)
	if dbResult.Error != nil {
		return dbResult.Error
	}

	return nil
}

func (db *DataBase) DeleteRecordByColumn(columnName string, value interface{}, row interface{}) error {
	// Выполняем запрос на удаление записи из таблицы, указанной в модели
	result := db.Where(fmt.Sprintf("%s = ?", columnName), value).Delete(row)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
