package main

import (
	"nstu/internal/logger"
	"nstu/internal/models"
	"nstu/internal/tg"
)

func main() {
	handleErr(logger.IniLogger())
	handleErr(models.InitConfig())
	handleErr(models.InitDataBase())

	bot, err := tg.InitTgBNot()
	handleErr(err)
	bot.HandleUpdates()
	tg.Bot = bot

}

// handleErr проверяет наличие ошибки. В случае наличия ошикби, логгирует ее и останавливает программу.
func handleErr(err error) {
	if err != nil {
		logger.Log.Error("Ошибка при запуске приложения: %s", err)
		panic("")
	}
}
