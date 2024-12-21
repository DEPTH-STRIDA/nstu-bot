package model

import (
	"time"

	"gorm.io/gorm"
)

// TimeSlot представляет временной слот (30 минут)
type TimeSlot struct {
	gorm.Model
	TeacherID   uint      `gorm:"not null"`
	StartTime   time.Time `gorm:"not null"`
	IsAvailable bool      `gorm:"not null"`
	StudentID   *uint
	Student     *Student `gorm:"foreignKey:StudentID"`
}

// Schedule представляет расписание преподавателя
type Schedule struct {
	gorm.Model
	TeacherID uint      `gorm:"not null;unique"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	WorkDays  string    `gorm:"type:varchar(7)"`
}

type DeletedAt struct {
	Time  time.Time `json:"time"`
	Valid bool      `json:"valid"` // Valid указывает, было ли установлено время удаления
}

// Subject представляет предмет
type Subject struct {
	gorm.Model
	Name string `gorm:"type:varchar(1000);not null;unique";json:"name"` // Имя предмета, уникальное
}

// LevelTraining представляет уровень подготовки
type LevelTraining struct {
	gorm.Model
	Name string `gorm:"type:varchar(1000);not null;unique";json:"name"` // Имя уровня подготовки, уникальное
}

type TutorExperience struct {
	gorm.Model
	Name  string `gorm:"type:varchar(1000);not null;unique";json:"name"`
	Value int    `json:"value"`
}

// Service представляет оказываемую услугу
type Service struct {
	gorm.Model
	Name string `gorm:"type:varchar(1000);not null;unique";json:"name"` // Название услуги
}

const (
	AdminRole   = "admin"
	TeacherRole = "teacher"
	StudentRole = "student"
)

// User представляет пользователя системы
type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(1000);not null";json:"name"`
	Email    string `gorm:"type:varchar(1000);not null;unique";json:"email"`
	Password string `gorm:"type:varchar(1000);not null";json:"password"`
	RoleID   uint   `gorm:"type:int;";json:"role_id"`
	Role     string `gorm:"type:varchar(100);not null";json:"role"`

	// Связи с профилями
	TeacherID *uint `json:"teacher_id,omitempty"`
	StudentID *uint `json:"student_id,omitempty"`
	AdminID   *uint `json:"admin_id,omitempty"`

	// Профили пользователя
	Teacher *Teacher `gorm:"foreignKey:TeacherID";json:"teacher,omitempty"`
	Student *Student `gorm:"foreignKey:StudentID";json:"student,omitempty"`
	Admin   *Admin   `gorm:"foreignKey:AdminID";json:"admin,omitempty"`
}

// Chat представляет чат между двумя пользователями
type Chat struct {
	gorm.Model
	User1ID uint `gorm:"not null"` // ID первого пользователя
	User2ID uint `gorm:"not null"` // ID второго пользователя
	User1   User `gorm:"foreignKey:User1ID"`
	User2   User `gorm:"foreignKey:User2ID"`
}

type ChatMessage struct {
	gorm.Model
	Name string `gorm:"type:varchar(1000);not null;unique"`
}

// Message представляет сообщение в чате
type Message struct {
	gorm.Model
	ChatID   uint   `gorm:"not null"` // ID чата
	SenderID uint   `gorm:"not null"` // ID отправителя
	Content  string `gorm:"type:text;not null"`
	Chat     Chat   `gorm:"foreignKey:ChatID"`
	Sender   User   `gorm:"foreignKey:SenderID"`
}

// ClassFormat представляет формат занятий
type ClassFormat struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null;unique"`
}

// Teacher представляет преподавателя
type Teacher struct {
	gorm.Model
	UserID        uint            `gorm:"not null"`
	Name          string          `gorm:"type:varchar(1000);not null"`
	Description   string          `gorm:"type:varchar(1000)"`
	Education     string          `gorm:"type:varchar(1000)"`
	ExperienceID  uint            `gorm:"not null"`
	Experience    TutorExperience `gorm:"foreignKey:ExperienceID"`
	Price         int             `gorm:"not null"`
	ClassFormats  []ClassFormat   `gorm:"many2many:teacher_class_formats"` // Изменено на many-to-many
	ImgUrl        string          `gorm:"type:varchar(1000)"`
	LevelTraining []LevelTraining `gorm:"many2many:teacher_level_trainings"`
	Subjects      []Subject       `gorm:"many2many:teacher_subjects"`
	Services      []Service       `gorm:"many2many:teacher_services"`
	Schedule      Schedule        `gorm:"foreignKey:TeacherID"`
	TimeSlots     []TimeSlot      `gorm:"foreignKey:TeacherID"`
}

// Admin представляет администратора
type Admin struct {
	gorm.Model
	Description string `gorm:"type:varchar(1000)";json:"description"` // Описание администратора
}

// Student представляет студента
type Student struct {
	gorm.Model
	Subjects    []Subject `gorm:"many2many:student_subjects";json:"subjects"` // Предметы
	Class       string    `gorm:"type:varchar(1000);not null"`                // Класс
	Description string    `gorm:"type:varchar(1000)"`                         // Описание администратора
}
