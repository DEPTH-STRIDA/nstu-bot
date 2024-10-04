package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Маршрутизатор
func (app *WebApp) SetRoutes() *mux.Router {
	router := mux.NewRouter()

	// Ограничение количества запросов от одного IP
	router.Use(LimitMiddleware)

	router.HandleFunc("/", app.HandleValidate).Methods("GET")
	router.HandleFunc("/main", InitDataValidation(app.HandleMain)).Methods("GET")

	router.HandleFunc("/get/groups", InitDataValidation(app.HandleGetGroups)).Methods("POST")
	router.HandleFunc("/get/my-groups", InitDataValidation(app.HandleGetMyGroups)).Methods("POST")
	router.HandleFunc("/get/group-schedule", InitDataValidation(app.HandleGetGroupSchedule)).Methods("POST")

	router.HandleFunc("/set/group-schedule", InitDataValidation(app.HandleSetGroupSchedule)).Methods("POST")

	router.HandleFunc("/group/join", InitDataValidation(app.HandleGroupJoin)).Methods("POST")
	router.HandleFunc("/group/exit", InitDataValidation(app.HandleExitGroup)).Methods("POST")

	router.HandleFunc("/group/exit-and-join", InitDataValidation(app.HandleGroupJoin)).Methods("POST")

	staticDir := "./ui/static/"
	fileServer := http.FileServer(http.Dir(staticDir))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	return router
}
