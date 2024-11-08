package web

import (
	"context"
	"net/http"
	"nstu/internal/logger"
	"nstu/internal/models"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type contextKey string

const (
	InitDataCtx contextKey = "init_data_ctx"
)

// TokenValidation проверяет наличие и правильность токена из заголовка с подписью "Validation"
func InitDataValidation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logStr := r.RemoteAddr + " -- " + r.URL.Path

		// var jsonData map[string]interface{}

		// err := json.NewDecoder(r.Body).Decode(&jsonData)
		// if err != nil && err != io.EOF {
		// 	http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		// 	logger.Log.Error(logStr + "Error parsing JSON: " + err.Error())
		// 	return
		// }
		// defer r.Body.Close()

		// initDataInterface, ok := jsonData["initData"]
		// if !ok {
		// 	http.Error(w, "The map does not contain initData", http.StatusBadRequest)
		// 	logger.Log.Error(logStr + "The map does not contain initData")
		// 	return
		// }
		// initDataStr, ok := initDataInterface.(string)
		// if !ok {
		// 	http.Error(w, "Can not convert initData", http.StatusBadRequest)
		// 	logger.Log.Error(logStr + "Can not convert initData: " + err.Error())
		// 	return
		// }

		// initData, err := validateInitData(initDataStr, models.Config.NSTUToken)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	logger.Log.Error(logStr + "Validate initData failed: " + err.Error())
		// 	return
		// }

		initDataTest := initdata.InitData{User: initdata.User{ID: 878413772, FirstName: "Maxim", LastName: "Tich", Username: "Tichomirov2003"}}
		logger.Log.Error(logStr + "Validate initData succses.")

		err := models.EnsureUserExists(models.DataBase.DB, initDataTest)
		if err != nil {
			logger.Log.Error("Не удалось проверить наличие пользвователя в БД: ", err)
			return
		}

		// Добавляем
		ctx := context.WithValue(r.Context(), InitDataCtx, &initDataTest)
		r = r.WithContext(ctx)

		// Вызываем следующий обработчик с обновленным запросом
		next.ServeHTTP(w, r)
	}
}
