package model

import "github.com/dgrijalva/jwt-go"

type AuthConfig struct {
	TimeToLive      int `envconfig:"AUTH_TIME_TO_LIVE" required:"true"`     // Время жизни кэша в минутах
	CleanupInterval int `envconfig:"AUTH_CLEANUP_INTERVAL" required:"true"` // Интервал очистки кэша в минутах
}

// Структура, которая кодируется в Json и передается вместе с HTTP.
type Token struct {
	UserId uint
	Role   string
	jwt.StandardClaims
}
