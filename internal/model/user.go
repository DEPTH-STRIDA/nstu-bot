package model

import "gorm.io/gorm"

type User struct {
	gorm.Model           // ID, CreatedAt, UpdatedAt, DeletedAt
	Username      string `gorm:"type:varchar(100)"`
	FirstName     string `gorm:"type:varchar(100)"`
	LastName      string `gorm:"type:varchar(100)"`
	Email         string `gorm:"type:varchar(500)"`
	Valid         bool   `gorm:"type:boolean"`
	MailExpiredAt string `gorm:"type:varchar(100)"`
	Code          string `gorm:"type:varchar(100)"`
	CRMID         int64  `gorm:"unique"`          // Уникальный индекс по ID из CRM
	TelegramID    int64  `gorm:"unique;not null"` // Уникальный индекс по ID из Telegram
}
