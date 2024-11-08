package models

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DataBase *DataBaseStruct

type DataBaseStruct struct {
	*gorm.DB
}

// InitDataBase подключается к БД по данным из конфига. Устанавливает подключение в глобальную переменную.
func InitDataBase() error {

	// Стандартная строка для подключения к БД postgresql
	dbUri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		Config.DataBaseConfig.UserName,
		Config.DataBaseConfig.Password,
		Config.DataBaseConfig.Host,
		Config.DataBaseConfig.Port,
		Config.DataBaseConfig.DataBaseName)

	// Сохранение "Подключения" в переменную
	conn, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		return err
	}
	DataBase = &DataBaseStruct{}

	DataBase.DB = conn

	return nil
}

type DatabaseService interface {
	// Возвращает пользователя по тг id
	GetUser(tgId int64) (*User, error)

	// Возвращает все существующие группы
	GetGroups() (*[]Group, error)
	// Возвращает группы определенного пользователя
	GetUserGroups(id uint) (*[]Group, error)
	// Возвращает определенную группу
	GetGroup(id uint) (*Group, error)

	// Обновляет все поля группы, кроме id
	UpdateGroup(id uint) error
	// Создает группу
	CreateGroup(group *Group) (*Group, error)

	// Устанавливает пользователю принадлежность к определенной группе
	JoinGroup(userId, groupId uint) error
	// Снимает принадлежность пользователя к какой-либо группе
	ExitGroup(id uint) error
}

// GetUser возвращает пользователя по tg id
func (db *DataBaseStruct) GetUser(tgId int64) (*User, error) {
	var user User
	result := db.DB.Where("tg_user_id = ?", tgId).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

// CreateUser создает пользователя в БД
func (db *DataBaseStruct) CreateUser(user *User) (*User, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}

	result := db.DB.Create(user)
	if result.Error != nil {
		// Проверяем, является ли ошибка нарушением уникальности tg_user_id
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, errors.New("user with this Telegram ID already exists")
		}
		return nil, result.Error
	}

	// Проверяем, был ли пользователь действительно создан
	if result.RowsAffected == 0 {
		return nil, errors.New("user was not created")
	}

	return user, nil
}

// CheckUserExists проверяет, существует ли пользователь с заданным Telegram ID в базе данных
func (db *DataBaseStruct) CheckUserExists(tgUserID uint) (bool, error) {
	var user User
	result := db.DB.Where("tg_user_id = ?", tgUserID).Limit(1).Find(&user)

	if result.Error != nil {
		// Если ошибка не связана с отсутствием записи, возвращаем её
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, result.Error
		}
		// Если запись не найдена, это не ошибка, а ожидаемый результат
		return false, nil
	}

	// Пользователь существует, если найдена хотя бы одна запись
	return result.RowsAffected > 0, nil
}

// Возвращает все существующие группы
func (db *DataBaseStruct) GetGroups() (*[]Group, error) {
	var groups []Group
	result := db.DB.Find(&groups)
	if result.Error != nil {
		return nil, result.Error
	}
	return &groups, nil
}

// Возвращает группы определенного пользователя
func (db *DataBaseStruct) GetUserGroups(id int64) (*[]Group, error) {
	var groups []Group
	result := db.DB.Where("holder_id = ?", id).Find(&groups)
	if result.Error != nil {
		return nil, result.Error
	}
	return &groups, nil
}

// Возвращает определенную группу
func (db *DataBaseStruct) GetGroup(id int64) (*Group, error) {
	var group Group
	result := db.DB.First(&group, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &group, nil
}

// Обновляет все поля группы, кроме id
func (db *DataBaseStruct) UpdateGroup(group *Group) error {
	if group == nil || group.ID == 0 {
		return errors.New("invalid group")
	}
	result := db.DB.Model(group).Updates(group)
	return result.Error
}

// Создает группу
func (db *DataBaseStruct) CreateGroup(group *Group) (*Group, error) {
	if group == nil {
		return nil, errors.New("group cannot be nil")
	}
	result := db.DB.Create(group)
	if result.Error != nil {
		return nil, result.Error
	}
	return group, nil
}

// Устанавливает пользователю принадлежность к определенной группе
func (db *DataBaseStruct) JoinGroup(userTgId int64, groupId uint) error {

	// Проверяем существование пользователя и группы
	user, err := db.GetUser(int64(userTgId))
	if err != nil {
		return errors.New("user not found: " + err.Error())
	}

	var group Group
	if err := db.DB.First(&group, groupId).Error; err != nil {
		return errors.New("group not found")
	}

	// Проверяем, не состоит ли пользователь уже в этой группе
	if user.ConsistsOf != nil && *user.ConsistsOf == int(groupId) {
		return errors.New("user already in this group")
	}

	// Обновляем поле ConsistsOf пользователя
	groupIdInt := int(groupId)
	user.ConsistsOf = &groupIdInt

	return db.DB.Save(&user).Error
}

// Снимает принадлежность пользователя к какой-либо группе
func (db *DataBaseStruct) ExitGroup(userTgId int64) error {
	// Проверяем существование пользователя и группы
	user, err := db.GetUser(userTgId)
	if err != nil {
		return errors.New("user not found: " + err.Error())
	}

	// Проверяем, состоит ли пользователь в какой-либо группе
	if user.ConsistsOf == nil || *user.ConsistsOf == 0 {
		return errors.New("user is not in any group")
	}

	// Устанавливаем ConsistsOf на 0, что означает отсутствие принадлежности к группе
	zeroValue := 0
	user.ConsistsOf = &zeroValue

	// Сохраняем изменения в базе данных
	if err := db.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
