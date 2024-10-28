package business

import (
	"app/internal/logger"
	"app/internal/model"
	"app/internal/moyklass"
	"app/internal/tg"
	"fmt"
	"log"
	"net/http"
)

var App *LogicApp

type LogicApp struct {
	mux     *http.ServeMux
	updates chan *moyklass.Update
}

func InitLogicApp(mux *http.ServeMux) error {
	App = &LogicApp{
		mux:     mux,
		updates: moyklass.MoyClass.GetUpdatesChan(),
	}

	return nil
}

func (app *LogicApp) StartServer() {
	// Создаем сервер
	server := &http.Server{
		Addr:    model.ConfigFile.APPIP + ":" + model.ConfigFile.APPPORT, // Порт для прослушивания
		Handler: app.mux,                                                 // Наш мультиплексор
	}
	logger.Log.Info("Starting server on ", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (app *LogicApp) StartHandleUpdates() {
	for v := range app.updates {
		fmt.Println("ПРИНЯЛИ СОБЫТИЕ: ", v)

		tg.TelegramBot.SendAllAdmins(fmt.Sprintf("ПРИНЯЛИ СОБЫТИЕ: %v", v.Object))
	}
}
