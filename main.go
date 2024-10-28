package main

import (
	"app/internal/business"
	"app/internal/db"
	"app/internal/logger"
	"app/internal/model"
	"app/internal/moyklass"
	"app/internal/tg"
	u "app/internal/utils"
	"net/http"
)

func main() {
	u.HandleFatalError(logger.IniLogger())

	u.HandleFatalError(model.InitConfig())

	u.HandleFatalError(db.InitDataBase())

	u.HandleFatalError(tg.InitTgBot())

	mux := http.NewServeMux()

	u.HandleFatalError(moyklass.InitNoyKlass(mux))

	u.HandleFatalError(business.InitLogicApp(mux))

	business.App.StartServer()
}
