package utils

import (
	"app/log"
	"app/model"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	mathrand "math/rand"
	"os"
	"regexp"
	"time"
	"unicode"

	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

// HandlerError логгирует ошибку в случае ее наличия и возвращает ее
func HandleError(err error) error {
	if err != nil {
		log.App.Error(err)
		return err
	}
	return nil
}

// HandleFatalError если err ошибка, то логгирует ее, отправляет всем админам в тг, если ошибки нет, то возвращает nil
func HandleFatalError(err error) error {
	if err != nil {
		log.App.Error("Критическая ошибка: ", err)

		os.Exit(1)
	}
	return nil
}

// setLocationTime устанавливает часовой пояс по умолчанию для глобальной переменной.
// Принимает строку с названием локации и возвращает ошибку, если часовой пояс не удалось загрузить.
func InitGlobalLocationTime() error {
	// Устанавливаем локацию по умолчанию для time.Local
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return fmt.Errorf("ошибка при смене локации на %s: %w", "Europe/Moscow", err)
	}
	time.Local = loc
	return nil
}

// HashPassword хеширует пароль с использованием bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.App.Error("Ошибка при хешировании пароля:", err)
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash проверяет соответствие пароля и хеша
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashEmailSHA256 хеширует email с использованием SHA-256
func HashEmailSHA256(email string) string {
	hash := sha256.Sum256([]byte(email))
	return hex.EncodeToString(hash[:])
}

// CompareEmailHash сравнивает хеши email
func CompareEmailHash(email, hash string) bool {
	return HashEmailSHA256(email) == hash
}

// ValidateEmail проверяет корректность адреса электронной почты
func ValidateEmail(email string) error {
	if len(email) > 999 {
		return errors.New("Адрес электронной почты слишком длинный (максимум 999 символов)")
	}

	re := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !re.MatchString(email) {
		return errors.New("Некорректный формат адреса электронной почты")
	}

	return nil
}

// ValidatePassword проверяет надежность пароля. Возвращает статус и ошибку в случае неудачи
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("Пароль должен содержать как минимум 8 символов")
	}
	if len(password) > 999 {
		return errors.New("Пароль слишком длинный (максимум 999 символа)")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("Пароль должен содержать хотя бы одну заглавную букву")
	}
	if !hasLower {
		return errors.New("Пароль должен содержать хотя бы одну строчную букву")
	}
	if !hasNumber {
		return errors.New("Пароль должен содержать хотя бы одну цифру")
	}
	if !hasSpecial {
		return errors.New("Пароль должен содержать хотя бы один специальный символ")
	}

	return nil
}

// GenerateCode генерирует случайный код из 5 цифр
func GenerateCode() string {
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano())) // Создаем локальный генератор случайных чисел
	code := r.Intn(100000)                                       // Генерация случайного числа от 0 до 99999
	return fmt.Sprintf("%05d", code)                             // Форматирование числа с ведущими нулями
}

// GenerateToken генерирует случайный токен заданной длины
func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes) // Используем crypto/rand
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// StructToJSONString принимает ссылку на структуру и возвращает её JSON-представление в виде строки.
// Если происходит ошибка, возвращает JSON-строку с описанием ошибки.
func StructToJSONString(v interface{}) string {
	// Преобразуем структуру в JSON
	bytes, err := json.Marshal(v)
	if err != nil {
		// Если произошла ошибка, возвращаем JSON-строку с описанием ошибки
		errorResponse := model.Response{
			Status:  "failure",
			Message: "Ошибка при преобразовании структуры в JSON: " + err.Error(),
		}
		errorBytes, _ := json.Marshal(errorResponse)
		return string(errorBytes)
	}
	// Преобразуем байты в строку и возвращаем
	return string(bytes)
}

// HashPasswordSHA256 хеширует пароль с использованием SHA-256
func HashPasswordSHA256(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
