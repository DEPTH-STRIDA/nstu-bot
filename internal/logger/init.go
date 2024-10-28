package logger

import (
	"fmt"
	"os"
)

// Log интерфейс логгера.
var Log Logger

// IniLogger иницилизация комбинированного логгера
func IniLogger() error {
	// Проверка существования директории
	if _, err := os.Stat("./log/"); os.IsNotExist(err) {
		return fmt.Errorf("директория /log/ не существует: %v", err)
	}

	logger, err := NewCombinedLogger("./log/")
	if err != nil {
		return fmt.Errorf("ошибка создания комбинированного логгера: %v", err)
	}

	Log = logger
	return nil
}
