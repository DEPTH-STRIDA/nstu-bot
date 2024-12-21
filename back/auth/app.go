// В данном пакеете реализуется бизнес-логика для авторизации и регистрации пользователей.
package auth

import (
	"app/cache"
	"app/config"
	"app/db"
	"app/log"
	"app/model"
	"app/smtp"
	"app/utils"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// CreateJWTToken создает jwt токен для данных пользователя
func CreateJWTToken(user *model.User) (string, error) {
	// Структура токена
	tk := &model.Token{
		UserId:         user.ID,
		Role:           user.Role,
		StandardClaims: jwt.StandardClaims{},
	}

	// Создание токена из структуры "tk" с алгоритмом (HMAC-SHA256) для шифрования токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)

	tokenString, err := token.SignedString([]byte(config.File.JWTTokenPassword))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWTToken разбирает JWT токен и возвращает данные п��льзователя
func ParseJWTToken(tokenString string) (*model.Token, error) {
	// Создаем экземпляр структуры Token для хранения данных из токена
	tk := &model.Token{}

	// Парсим токен с использованием секретного ключа
	token, err := jwt.ParseWithClaims(tokenString, tk, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи токена соответствует ожидаемому
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.File.JWTTokenPassword), nil
	})

	if err != nil {
		return nil, err
	}

	// Проверяем, является ли токен действительным и содержит ли он ожидаемые данные
	if claims, ok := token.Claims.(*model.Token); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// CheckEmailAvailability проверяет, занята ли почта.
// Если true, то почта занята.
func CheckEmailAvailability(email string) error {
	temp := &model.User{}

	// Проверка на дубликаты почты
	err := db.App.Where("email = ?", email).First(temp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return fmt.Errorf("Почта занята")
}

// Validate проверяет корректность пароля и почты
func Validate(user *model.User) error {
	// Проверка электронной почты
	err := utils.ValidateEmail(user.Email)
	if err != nil {
		return err
	}

	// Проверка пароля
	err = utils.ValidatePassword(user.Password)
	if err != nil {
		return err
	}

	// Проверка на занятость почты
	err = CheckEmailAvailability(user.Email)
	if err != nil {
		return err
	}

	return nil
}

// RegisterUser регистрирует пользователя. Проверяет корректность почты и пароля, отправляет на почту письмо с кодом подтверждения.
func RegisterStudent(email, password string) (*model.RoleResponse, string, error) {
	user := model.User{
		Email:    email,
		Password: password,
	}

	err := Validate(&user)
	if err != nil {
		return nil, "", err
	}

	code := utils.GenerateCode()

	err = smtp.App.SendConfirmationCodeEmail(user.Email, code, smtp.RegistrationCode)
	if err != nil {
		log.App.Error(" failed to send registration code: ", err)
		return nil, "", err
	}

	token, err := utils.GenerateToken(32)
	if err != nil {
		return nil, "", err
	}

	cache.Auth.Set(token, model.CachedUser{
		Email:      user.Email,
		Code:       code,
		Password:   user.Password,
		ActionType: model.RegistrationStarted,
	})

	return &model.RoleResponse{
		Response: model.Response{
			Status:  "success",
			Message: "Регистрация начата. На почту отправлен код подтверждения.",
		},
	}, token, nil
}

// ValidateEmail проверяет доступность email для регистрации.
func ValidateEmail(email string) error {
	// Проверка электронной почты
	err := utils.ValidateEmail(email)
	if err != nil {
		return err
	}

	// Проверка на занятость почты
	err = CheckEmailAvailability(email)
	if err != nil {
		return err
	}

	return nil
}

// ConfirmRegistration подтверждает регистрацию пользователя.
func ConfirmRegistration(jwtToken, code string) (*model.RoleResponse, string, error) {
	user, ok := cache.Auth.Get(jwtToken)
	if !ok {
		return nil, "", fmt.Errorf("Данный аккаунт не требует подтверждения")
	}

	if user.Code != code {
		return nil, "", fmt.Errorf("Неверный код подтверждения")
	}

	if user.ActionType != model.RegistrationStarted {
		return nil, "", fmt.Errorf("ошибка при подтверждении регистрации")
	}

	userPostfix, err := utils.GenerateToken(5)
	if err != nil {
		return nil, "", err
	}

	newUser := &model.User{
		Email:    user.Email,
		Password: user.Password,
		Role:     model.StudentRole,
		Name:     "User" + userPostfix,
	}

	err = Validate(newUser)
	if err != nil {
		return nil, "", err
	}
	newPassword := utils.HashPasswordSHA256(user.Password)

	newUser.Password = newPassword

	res := db.App.Create(newUser)
	if res.Error != nil {
		return nil, "", res.Error
	}

	newStudent := &model.Student{}

	res = db.App.Create(newStudent)
	if res.Error != nil {
		return nil, "", res.Error
	}

	res = db.App.Model(&model.User{}).Where("email = ?", user.Email).Updates(model.User{RoleID: newStudent.ID})
	if res.Error != nil {
		return nil, "", res.Error
	}

	token, err := CreateJWTToken(newUser)
	if err != nil {
		return nil, "", err
	}

	return &model.RoleResponse{
		Role: model.StudentRole,
		Name: newUser.Name,
		Response: model.Response{
			Status:  "success",
			Message: "Регистрация завершена",
		},
	}, token, nil
}

// Login выполняет авторизацию пользователя. Возвращает токен и ошибку в случае неудачи.
func Login(email, password string) (*model.RoleResponse, string, error) {
	// Проверка, существует ли пользователь с данным email
	var account model.User
	err := db.App.Where("email = ?", email).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", fmt.Errorf("Пользователь с такой почтой не найден")
		}
		return nil, "", err
	}
	newPassword := utils.HashPasswordSHA256(password)
	if newPassword != account.Password {
		return nil, "", fmt.Errorf("Неверный пароль")
	}

	token, err := CreateJWTToken(&account)
	if err != nil {
		return nil, "", err
	}

	return &model.RoleResponse{
		Role: account.Role,
		Response: model.Response{
			Status:  "success",
			Message: "Авторизация прошла успешно",
		},
	}, token, nil
}

func JwtLogin(tokenString string) (*model.RoleResponse, error) {
	// Извлечение токена из заголовка
	token, err := ParseJWTToken(tokenString)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = db.App.First(&user, token.UserId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Пользователь не найден")
		}
		return nil, err
	}

	if user.Role != token.Role {
		return nil, fmt.Errorf("Несоответствие роли пользователя")
	}

	return &model.RoleResponse{
		Role: user.Role,
		Name: user.Name,
		Response: model.Response{
			Status:  "success",
			Message: "Вход выполнен",
		},
	}, nil
}

func ResetPassword(email string) (string, error) {
	err := utils.ValidateEmail(email)
	if err != nil {
		return "", err
	}

	code := utils.GenerateCode()

	err = smtp.App.SendConfirmationCodeEmail(email, code, smtp.PasswordResetCode)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(32)
	if err != nil {
		return "", err
	}

	cache.Auth.Set(token, model.CachedUser{
		Email:      email,
		Code:       code,
		ActionType: model.PasswordChangeStarted,
	})

	return token, nil
}

func ConfirmResetPassword(jwtToken, code string) (string, error) {
	user, ok := cache.Auth.Get(jwtToken)
	if !ok {
		return "", fmt.Errorf("Данный аккаунт не требует подтверждения")
	}

	if user.Code != code {
		return "", fmt.Errorf("Неверный код подтверждения")
	}

	if user.ActionType != model.RegistrationStarted {
		return "", fmt.Errorf("Ошибка при подтверждении регистрации")
	}

	newUser := &model.CachedUser{
		Email:      user.Email,
		ActionType: model.PasswordChangeComplete,
	}

	cache.Auth.Delete(jwtToken)

	newToken, err := utils.GenerateToken(32)
	if err != nil {
		return "", err
	}

	cache.Auth.Set(newToken, *newUser)

	return newToken, nil
}

func SetNewPassword(jwtToken, password string) (string, error) {
	user, ok := cache.Auth.Get(jwtToken)
	if !ok {
		return "", fmt.Errorf("Данный аккаунт не требует подтверждения")
	}

	if user.ActionType != model.PasswordChangeComplete {
		return "", fmt.Errorf("ошибка при установке нового пароля")
	}

	err := utils.ValidatePassword(password)
	if err != nil {
		return "", err
	}

	newPassword := utils.HashPasswordSHA256(user.Password)

	res := db.App.Model(&model.User{}).Where("email = ?", user.Email).Updates(model.User{Password: newPassword})
	if res.Error != nil {
		return "", res.Error
	}

	newUser := &model.User{}
	err = db.App.Where("email = ?", user.Email).First(newUser).Error
	if err != nil {
		return "", err
	}

	token, err := CreateJWTToken(newUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

func UpdateUser(jwtToken string, req model.UpdateUserRequest) error {
	// Извлечение токена из заголовка
	token, err := ParseJWTToken(jwtToken)
	if err != nil {
		return err
	}

	var user model.User
	err = db.App.First(&user, token.UserId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Пользователь не найден")
		}
		return err
	}

	if user.Role != token.Role {
		return fmt.Errorf("Несоответствие роли пользователя")
	}

	res := &gorm.DB{}

	switch user.Role {
	case model.StudentRole:
		res = db.App.Model(&model.Student{}).Where("id = ?", user.ID).Updates(req.Student)
	case model.TeacherRole:
		res = db.App.Model(&model.Teacher{}).Where("id = ?", user.ID).Updates(req.Teacher)
	case model.AdminRole:
		res = db.App.Model(&model.Admin{}).Where("id = ?", user.ID).Updates(req.Admin)
	}

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func GetProfileInfo(jwtToken string) (*model.Admin, *model.Teacher, *model.Student, error) {
	// Парсим токен
	token, err := ParseJWTToken(jwtToken)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Получаем пользователя с предзагрузкой связанных данных
	var user model.User
	if err := db.App.
		Preload("Admin").
		Preload("Teacher.Subjects").
		Preload("Teacher.LevelTraining").
		Preload("Teacher.Experience").
		Preload("Teacher.Services").
		Preload("Student.Subjects").
		First(&user, token.UserId).Error; err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Возвращаем соответствующий профиль
	switch user.Role {
	case model.AdminRole:
		if user.Admin == nil {
			return nil, nil, nil, fmt.Errorf("admin profile not found")
		}
		return user.Admin, nil, nil, nil
	case model.TeacherRole:
		if user.Teacher == nil {
			return nil, nil, nil, fmt.Errorf("teacher profile not found")
		}
		return nil, user.Teacher, nil, nil
	case model.StudentRole:
		if user.Student == nil {
			return nil, nil, nil, fmt.Errorf("student profile not found")
		}
		return nil, nil, user.Student, nil
	default:
		return nil, nil, nil, fmt.Errorf("unknown user role: %s", user.Role)
	}
}

func SaveUserInfo(jwtToken string, req model.UpdateUserRequest) error {
	// Парсим токен
	token, err := ParseJWTToken(jwtToken)
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	// Получаем пользователя
	var user model.User
	if err := db.App.First(&user, token.UserId).Error; err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Обновляем данные в зависимости от роли
	switch user.Role {
	case model.AdminRole:
		if user.AdminID == nil {
			return fmt.Errorf("admin profile not found")
		}
		// Обновляем только непустые поля админа
		updates := map[string]interface{}{}
		if req.Admin.Description != "" {
			updates["description"] = req.Admin.Description
		}
		if len(updates) > 0 {
			if err := db.App.Model(&model.Admin{}).Where("id = ?", *user.AdminID).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update admin: %w", err)
			}
		}

	case model.TeacherRole:
		if user.TeacherID == nil {
			return fmt.Errorf("teacher profile not found")
		}
		// Обновляем только непустые поля учителя
		updates := map[string]interface{}{}
		if req.Teacher.Name != "" {
			updates["name"] = req.Teacher.Name
		}
		if req.Teacher.Description != "" {
			updates["description"] = req.Teacher.Description
		}
		if req.Teacher.Education != "" {
			updates["education"] = req.Teacher.Education
		}
		if req.Teacher.ExperienceID != 0 {
			updates["experience_id"] = req.Teacher.ExperienceID
		}
		if req.Teacher.Price != 0 {
			updates["price"] = req.Teacher.Price
		}
		// if req.Teacher.ClassFormat != "" {
		// 	updates["class_format"] = req.Teacher.ClassFormat
		// }
		if req.Teacher.ImgUrl != "" {
			updates["img_url"] = req.Teacher.ImgUrl
		}
		if len(updates) > 0 {
			if err := db.App.Model(&model.Teacher{}).Where("id = ?", *user.TeacherID).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update teacher: %w", err)
			}
		}
		// Обновляем связи, если они предоставлены
		if len(req.Teacher.Subjects) > 0 {
			if err := db.App.Model(&model.Teacher{}).Where("id = ?", *user.TeacherID).Association("Subjects").Replace(req.Teacher.Subjects); err != nil {
				return fmt.Errorf("failed to update teacher subjects: %w", err)
			}
		}
		if len(req.Teacher.LevelTraining) > 0 {
			if err := db.App.Model(&model.Teacher{}).Where("id = ?", *user.TeacherID).Association("LevelTraining").Replace(req.Teacher.LevelTraining); err != nil {
				return fmt.Errorf("failed to update teacher levels: %w", err)
			}
		}
		if len(req.Teacher.Services) > 0 {
			if err := db.App.Model(&model.Teacher{}).Where("id = ?", *user.TeacherID).Association("Services").Replace(req.Teacher.Services); err != nil {
				return fmt.Errorf("failed to update teacher services: %w", err)
			}
		}

	case model.StudentRole:
		if user.StudentID == nil {
			return fmt.Errorf("student profile not found")
		}
		// Обновляем только непустые поля студента
		updates := map[string]interface{}{}
		if req.Student.Class != "" {
			updates["class"] = req.Student.Class
		}
		if req.Student.Description != "" {
			updates["description"] = req.Student.Description
		}
		if len(updates) > 0 {
			if err := db.App.Model(&model.Student{}).Where("id = ?", *user.StudentID).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update student: %w", err)
			}
		}
		// Обновляем предметы студента, если они предоставлены
		if len(req.Student.Subjects) > 0 {
			if err := db.App.Model(&model.Student{}).Where("id = ?", *user.StudentID).Association("Subjects").Replace(req.Student.Subjects); err != nil {
				return fmt.Errorf("failed to update student subjects: %w", err)
			}
		}

	default:
		return fmt.Errorf("unknown user role: %s", user.Role)
	}

	return nil
}
