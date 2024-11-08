package main

import (
	"nstu/internal/logger"
	"nstu/internal/models"
	"nstu/internal/tg"
	"nstu/internal/web"
)

func main() {
	handleErr(logger.IniLogger())
	handleErr(models.InitConfig())
	handleErr(models.InitDataBase())

	bot, err := tg.InitTgBNot()
	handleErr(err)
	go bot.HandleUpdates()
	tg.Bot = bot

	webApp, err := web.NewWebApp()
	handleErr(err)

	handleErr(webApp.HandleUpdates())
}

// handleErr проверяет наличие ошибки. В случае наличия ошикби, логгирует ее и останавливает программу.
func handleErr(err error) {
	if err != nil {
		logger.Log.Error("Ошибка при запуске приложения: %s", err)
		panic("")
	}
}
