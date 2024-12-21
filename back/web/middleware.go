package web

import (
	"app/auth"
	"app/model"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// ContextKey определяет тип для ключей в контексте
type ContextKey string

const (
	// UserIDKey используется для хранения ID пользователя в контексте
	UserIDKey ContextKey = "userID"
	// UserRoleKey используется для хранения роли пользователя в контексте
	UserRoleKey ContextKey = "userRole"
)

// JWTMiddleware проверяет JWT токен и добавляет данные о пользователе в контекст
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response := model.Response{Status: "error", Message: "Authorization header missing"}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Извлечение токена из заголовка
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := auth.ParseJWTToken(tokenString)
		if err != nil {
			response := model.Response{Status: "error", Message: "Invalid token"}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Добавление данных о пользователе в контекст
		ctx := context.WithValue(r.Context(), UserIDKey, token.UserId)
		ctx = context.WithValue(ctx, UserRoleKey, token.Role)

		// Передача управления следующему обработчику
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
