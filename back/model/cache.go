package model

type ActionType int

const (
	RegistrationStarted    ActionType = iota // Начало регистрации
	RegistrationComplete                     // Завершение регистрации
	Login                                    // Вход в аккаунт
	PasswordChangeStarted                    // Начало смены пароля
	PasswordChangeComplete                   // Завершение смены пароля
)

// PerformAction выполняет действие в зависимости от типа действия пользователя
func PerformAction(user *CachedUser) string {
	switch user.ActionType {
	case RegistrationComplete:
		return "Завершение регистрации"
	case Login:
		return "Вход в аккаунт"
	case PasswordChangeComplete:
		return "Завершение смены пароля"
	default:
		return "Неизвестное действие"
	}
}

type CachedUser struct {
	Name       string
	Email      string
	Code       string
	Password   string
	ActionType ActionType
}
