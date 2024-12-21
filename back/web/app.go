package web

import (
	"app/config"
	"app/log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// WebSocketServer представляет сервер WebSocket
type WebApp struct {
	Router   *mux.Router // Маршрутизатор
	Upgrader websocket.Upgrader
}

// NewWebSocketServer создает новый сервер WebSocket
func NewWebSocketServer() *WebApp {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	app := &WebApp{
		Router:   mux.NewRouter(),
		Upgrader: upgrader,
	}
	app.SetRoutes()

	return app
}

// StartServer запускает сервер. Данные для запуска берутся из конфига.
func (app *WebApp) StartServer() error {
	conf := config.File.WebConfig

	log.App.Info("Start server: ", conf.APPIP+":"+conf.APPPORT)
	err := http.ListenAndServe(conf.APPIP+":"+conf.APPPORT, app.Router)
	if err != nil {
		log.App.Info("ListenAndServe: ", err)
		return err
	}
	return nil
}
