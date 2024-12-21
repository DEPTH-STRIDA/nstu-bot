package main

import (
	"app/cache"
	"app/config"
	"app/db"
	_ "app/docs" // Не удалять. Для SWAGGER!
	"app/log"
	"app/smtp"
	u "app/utils"
	"app/web"
)

// @title Prosto Matika API
// @version 1.0
// @description API для управления данными и взаимодействия через HTTP-эндпоинты.
// Этот API предоставляет набор функций для работы с данными, обеспечивая надежную и быструю обработку запросов.

// @BasePath /api/v1
func main() {

	u.HandleFatalError(log.Init())

	u.HandleFatalError(u.InitGlobalLocationTime())

	u.HandleFatalError(config.Init())

	u.HandleFatalError(cache.Init())

	u.HandleFatalError(smtp.Init())

	u.HandleFatalError(db.Init())

	u.HandleFatalError(web.Init())

	u.HandleFatalError(web.App.StartServer())
}
