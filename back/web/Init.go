package web

var App *WebApp

func Init() error {
	App = NewWebSocketServer()
	return nil
}
