package models

import (
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"gorm.io/gorm"
)

// User представляет таблицу users
type User struct {
	gorm.Model
	TgUserID   int64  `gorm:"column:tg_user_id;uniqueIndex;not null"`
	ConsistsOf *int   `gorm:"column:consists_of"`
	FirstName  string `gorm:"column:first_name"`
	SecondName string `gorm:"column:second_name"`
	UserName   string `gorm:"column:user_name"`
}

// TableName определяет имя таблицы для User
func (User) TableName() string {
	return "users"
}
func EnsureUserExists(db *gorm.DB, initData initdata.InitData) error {
	var user User

	// Check if user exists in the database
	result := db.Where("tg_user_id = ?", initData.User.ID).First(&user)

	if result.Error == nil {
		// User already exists, no action needed
		return nil
	}

	if result.Error != gorm.ErrRecordNotFound {
		// An unexpected error occurred
		return result.Error
	}

	// User doesn't exist, create a new record
	newUser := User{
		TgUserID:   initData.User.ID,
		FirstName:  initData.User.FirstName,
		SecondName: initData.User.LastName,
		UserName:   initData.User.Username,
	}

	result = db.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
