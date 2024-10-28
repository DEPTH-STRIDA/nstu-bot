package db

import (
	"app/internal/logger"
	"app/internal/model"
	"app/internal/utils"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *DataBase

type DataBase struct {
	*gorm.DB
}

// NewDatabase создает новое подключение к базе данных
func InitDataBase() error {
	conf := model.ConfigFile.DataBaseConfig

	dsn := ""
	if conf.Host != "" {
		dsn += "host=" + conf.Host + "\n"
	}
	if conf.UserName != "" {
		dsn += "user=" + conf.UserName + "\n"
	}
	if conf.DBName != "" {
		dsn += "dbname=" + conf.DBName + "\n"
	}
	if conf.Password != "" {
		dsn += "password=" + conf.Password + "\n"
	}
	if conf.SSLMode != "" {
		dsn += "sslmode=" + conf.SSLMode + "\n"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	DB = &DataBase{
		DB: db,
	}

	DB.AutoMigrate(&model.User{})

	logger.Log.Info("Подключение к БД успешно установлено.")
	return nil
}

// CheckCode проверяет, что присланный код соответствует отправленному, проверяет время, делает аккаунт валидным
func (db *DataBase) CheckCode(telegramID, code string) error {
	var user model.User
	res := db.DB.Where("telegram_id = ?", telegramID).First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			logger.Log.Error("запись не найдена, сперва необходимо начать авторизацию: ", res.Error)
			return fmt.Errorf("запись не найдена, сперва необходимо начать авторизацию")
		}
		logger.Log.Error("ошибка при работе с БД: ", res.Error)
		return fmt.Errorf("ошибка при работе с БД: %w", res.Error)
	}

	// Почта истекает в
	mailExiredAt, err := utils.StringToTime(user.MailExpiredAt)
	if err != nil {
		logger.Log.Error("ошибка при получении времени истекания почты: ", err)
		return fmt.Errorf("ошибка при работе со временем: %w", res.Error)
	}

	if time.Now().After(mailExiredAt) {
		logger.Log.Error("срок годности кода подтверждения истек. ")
		return fmt.Errorf(`срок действия кода истек. Нажмите "Отправить код повторно" для получения нового кода подтверждения`)
	}

	oldCode, err := utils.DcryptToken(user.Code, model.ConfigFile.MailPassword)
	if err != nil {
		logger.Log.Error("ошибка при шифровании нового кода: ", err)
		return fmt.Errorf("%w", res.Error)
	}

	if oldCode != code {
		logger.Log.Error("неверный код")
		return fmt.Errorf("неверный код")
	}

	// Обновить одно поле
	db.Model(&model.User{}).
		Where("telegram_id = ?", telegramID).
		Update("valid", true)
	return nil
}

// SetNewMail устанавливает шифрованную почту в БД аккаунту, сбрасывает авторизованность, коды
func (db *DataBase) SetNewMail(telegramID, mail string) error {
	newMail, err := utils.EncryptToken(mail, model.ConfigFile.MailPassword)
	if err != nil {
		return err
	}

	resp := db.Model(&model.User{}).
		Where("telegram_id = ?", telegramID).
		Updates(map[string]interface{}{
			"code":            "",
			"email":           newMail,
			"valid":           false,
			"mail_expired_at": utils.TimeToString(time.Now().Add(time.Hour)),
		})

	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

// SetNewCode шифрует новый токен, устанавливает его и временную метку истекания аккаунту
func (db *DataBase) SetNewCode(telegramID, code string) error {
	newCode, err := utils.EncryptToken(code, model.ConfigFile.MailPassword)
	if err != nil {
		return err
	}

	resp := db.Model(&model.User{}).
		Where("telegram_id = ?", telegramID).
		Updates(map[string]interface{}{
			"code":            newCode,
			"mail_expired_at": utils.TimeToString(time.Now().Add(time.Hour)),
		})
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

// UserIsValid проверяет авторизован ли пользователь
func (db *DataBase) UserIsValid(telegramID int64) (bool, error) {
	var user model.User
	res := db.DB.Where("telegram_id = ?", telegramID).First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			logger.Log.Error("запись не найдена, сперва необходимо начать авторизацию: ", res.Error)
			return false, fmt.Errorf("запись не найдена, сперва необходимо начать авторизацию")
		}
		logger.Log.Error("ошибка при работе с БД: ", res.Error)
		return false, fmt.Errorf("ошибка при работе с БД: %w", res.Error)
	}

	if user.Valid {
		return true, nil
	}

	return false, nil
}

func (db *DataBase) CreateUserIfNotExistsWithConditions(newUser model.User) {
	var user model.User

	result := db.Where("telegram_id = ?", newUser.TelegramID).First(&user)

	if result.RowsAffected == 0 {
		if err := db.Create(&newUser).Error; err != nil {
			return
		}
	}
}

func (db *DataBase) Exit(telegramID int64) error {
	resp := db.Model(&model.User{}).
		Where("telegram_id = ?", telegramID).
		Updates(map[string]interface{}{
			"code":  "",
			"email": "",
			"valid": false,
		})

	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

func (db *DataBase) GetMail(telegramID int64) (string, error, bool) {
	var user model.User
	res := db.DB.Where("telegram_id = ?", telegramID).First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			logger.Log.Error("запись не найдена, сперва необходимо начать авторизацию: ", res.Error)
			return "", fmt.Errorf("запись не найдена, сперва необходимо начать авторизацию"), false
		}
		logger.Log.Error("ошибка при работе с БД: ", res.Error)
		return "", fmt.Errorf("ошибка при работе с БД: %w", res.Error), false
	}

	newMail, err := utils.DcryptToken(user.Email, model.ConfigFile.MailPassword)
	if err != nil {
		return "", err, false
	}

	// Почта истекает в
	mailExiredAt, err := utils.StringToTime(user.MailExpiredAt)
	if err != nil {
		logger.Log.Error("ошибка при получении времени истекания почты: ", err)
		return "", fmt.Errorf("ошибка при работе со временем: %w", res.Error), false
	}

	// Проверка времени
	if time.Now().After(mailExiredAt) {
		logger.Log.Error("срок годности кода подтверждения истек. ")
		return "", fmt.Errorf(`срок действия кода истек. Нажмите "Отправить код повторно" для получения нового кода подтверждения`), false
	}

	return newMail, nil, true
}
