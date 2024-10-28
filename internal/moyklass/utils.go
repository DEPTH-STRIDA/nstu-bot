package moyklass

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidateRoutePath проверяет и форматирует путь маршрута
// Возвращает отформатированный путь и ошибку, если путь невалидный
func validateRoutePath(path string) (string, error) {
	// Проверка на пустую строку
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Удаляем пробелы в начале и конце
	path = strings.TrimSpace(path)

	// Проверка на недопустимые символы
	for _, r := range path {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '/' && r != '-' && r != '_' {
			return "", fmt.Errorf("path contains invalid character: %c", r)
		}
	}

	// Добавляем / в начало, если его нет
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Удаляем / в конце, если он есть
	path = strings.TrimSuffix(path, "/")

	return path, nil
}
