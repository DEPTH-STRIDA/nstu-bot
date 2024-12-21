package web

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	basePath = "/api/v1"
)

// NotFoundHandler - обработчик для несуществующих маршрутов.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	log.Printf("404 Not Found: %s", r.URL.Path)
}

// SetRoutes устанавливает маршруты для HTTP-сервера.
// Эта функция связывает URL-пути с обработчиками.
func (app *WebApp) SetRoutes() {
	// Маршрут для WebSocket соединения
	app.Router.Handle(basePath+"/ws", JWTAuth(http.HandlerFunc(app.HandleConnections)))

	// Маршрут для Swagger UI
	app.Router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods("GET")

	// Маршруты для аутентификации и регистрации
	app.Router.HandleFunc(basePath+"/auth/login", app.HandleLogin).Methods("POST")
	app.Router.HandleFunc(basePath+"/auth/jwt", app.HandleJwtLogin).Methods("POST")
	app.Router.HandleFunc(basePath+"/registration/start", app.HandleRegistrationStarted).Methods("POST")
	app.Router.HandleFunc(basePath+"/registration/confirm", app.HandleRegistrationConfirmation).Methods("POST")
	app.Router.HandleFunc(basePath+"/registration/validate-email", app.HandleValidateEmail).Methods("POST")

	app.Router.HandleFunc(basePath+"/ws/teachers", app.HandleTeacherSearch).Methods("GET")

	app.Router.HandleFunc(basePath+"/user/update", app.HandleUpdateUser).Methods("POST")

	app.Router.HandleFunc(basePath+"/user/get-search-data", app.HandleGetSearchData).Methods("POST")
	app.Router.HandleFunc(basePath+"/user/get-profile-info", app.HandleGetProfileInfo).Methods("GET")
	app.Router.HandleFunc(basePath+"/user/save-user-info", app.HandleSaveUserInfo).Methods("POST")
	app.Router.HandleFunc(basePath+"/user/get-search-info", app.HandleGetSearchData).Methods("POST")

	// Маршруты для выхода
	app.Router.HandleFunc(basePath+"/auth/logout", app.Logout).Methods("POST")

	// Маршруты для сброса пароля
	app.Router.HandleFunc(basePath+"/password/reset/start", app.HandleResetPasswordStarted).Methods("POST")
	app.Router.HandleFunc(basePath+"/password/reset/confirm", app.HandleResetPasswordConfirmation).Methods("POST")
	app.Router.HandleFunc(basePath+"/password/new", app.HandleSetNewPassword).Methods("POST")

	// Маршруты для аутентификации по почте
	app.Router.HandleFunc(basePath+"/mail/auth", app.HandleMailAuthentication).Methods("POST")
	app.Router.HandleFunc(basePath+"/mail/auth/confirm", app.HandleMailAuthConfirmation).Methods("POST")

	// Обработчик для несуществующих маршрутов
	app.Router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
}
