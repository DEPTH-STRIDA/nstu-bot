package web

import (
	"fmt"
	"net/http"
	"nstu/internal/logger"
	"nstu/internal/models"
	"nstu/internal/request"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

// WebApp веб приложение. Мозг программы, который использует большинство других приложений
type WebApp struct {
	Router        *mux.Router                   // Маршрутизатор
	TemplateCache map[string]*template.Template // Карта шаблонов

	reqSheet    *request.RequestHandler
	reqTelegram *request.RequestHandler
}

// NewWebApp создает и возвращает веб приложение
func NewWebApp() (*WebApp, error) {
	// Загрузка шаблонов
	templateCache, err := NewTemplateCache("./ui/html/")
	if err != nil {
		return nil, err
	}

	reqSheet, err := request.NewRequestHandler(100)
	if err != nil {
		return nil, err
	}
	go reqSheet.ProcessRequests(1 * time.Second)

	reqTelegram, err := request.NewRequestHandler(100)
	if err != nil {
		return nil, err
	}
	go reqTelegram.ProcessRequests(1 * time.Second)

	app := WebApp{
		TemplateCache: templateCache,
		reqSheet:      reqSheet,
		reqTelegram:   reqTelegram,
	}
	// Установка параметров
	app.Router = app.SetRoutes()
	return &app, nil
}

// HandleUpdates запускает HTTP сервер
func (app *WebApp) HandleUpdates() error {
	logger.Log.Info("Запуск сервера по адрессу " + models.Config.WebAppConfig.AppIP + ":" + models.Config.WebAppConfig.AppPort)

	err := http.ListenAndServe(models.Config.WebAppConfig.AppIP+":"+models.Config.WebAppConfig.AppPort, app.Router)
	if err != nil {
		return fmt.Errorf("ошибка при запуске сервера: %v", err)
	}
	return nil
}

func validateInitData(initDataStr, token string) (*initdata.InitData, error) {
	expIn := 1 * time.Hour
	err := initdata.Validate(initDataStr, token, expIn)
	if err != nil {
		return nil, err
	}
	initData, err := initdata.Parse(initDataStr)
	if err != nil {
		return nil, err
	}
	return &initData, nil
}
