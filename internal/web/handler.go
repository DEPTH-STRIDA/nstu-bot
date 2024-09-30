package web

import (
	"encoding/json"
	"io"
	"net/http"
	"nstu/internal/logger"
	"nstu/internal/models"
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

func (wa *WebApp) HandleValidate(w http.ResponseWriter, r *http.Request) {
	err := wa.render(w, "validate.page.tmpl", nil)
	if err != nil {
		logger.Log.Error("Не удалось выполнить рендер: " + err.Error())
		http.Error(w, "Не удалось выполнить рендер: "+err.Error(), http.StatusInternalServerError)
	}
}

func (wa *WebApp) HandleMain(w http.ResponseWriter, r *http.Request) {
	err := wa.render(w, "main.page.tmpl", nil)
	if err != nil {
		logger.Log.Error("Не удалось выполнить рендер: " + err.Error())
		http.Error(w, "Не удалось выполнить рендер: "+err.Error(), http.StatusInternalServerError)
	}
}

func (wa *WebApp) HandleGetGroups(w http.ResponseWriter, r *http.Request) {

	initData, ok := r.Context().Value(InitDataCtx).(*initdata.InitData)
	if !ok || initData == nil {
		logger.Log.Error("Не удалось извлечь initData из контекста")
		http.Error(w, "Не удалось извлечь initData из контекста", http.StatusInternalServerError)
		return
	}

	groups, err := models.DataBase.GetGroups()
	if err != nil {
		logger.Log.Error("Не удалось получить группы из БД: ", err)
		http.Error(w, "Не удалось получить группы из БД: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := models.DataBase.GetUser(initData.User.ID)
	if err != nil {
		logger.Log.Error("Не удалось получить пользователя из БД: ", err)
		http.Error(w, "Не удалось получить пользователя из БД: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseBody := map[string]interface{}{"groups": groups, "consists-of": user.ConsistsOf}

	responseBodyJSON, err := json.Marshal(responseBody)
	if err != nil {
		logger.Log.Error("Не удалось выполнить маршалинг JSON: ", err)
		http.Error(w, "Не удалось выполнить маршалинг JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Запишите JSON в тело ответа
	_, err = w.Write(responseBodyJSON)
	if err != nil {
		logger.Log.Error("Не удалось записать JSON в тело ответа: ", err)
		// Здесь мы уже начали отправлять ответ, поэтому можем только залогировать ошибку
		return
	}
}

func (wa *WebApp) HandleGetMyGroups(w http.ResponseWriter, r *http.Request) {

	initData, ok := r.Context().Value(InitDataCtx).(*initdata.InitData)
	if !ok || initData == nil {
		logger.Log.Error("Не удалось извлечь initData из контекста")
		http.Error(w, "Не удалось извлечь initData из контекста", http.StatusInternalServerError)
		return
	}

	groups, err := models.DataBase.GetUserGroups(initData.User.ID)
	if err != nil {
		logger.Log.Error("Не удалось получить группы из БД: ", err)
		http.Error(w, "Не удалось получить группы из БД: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := models.DataBase.GetUser(initData.User.ID)
	if err != nil {
		logger.Log.Error("Не удалось получить пользователя из БД: ", err)
		http.Error(w, "Не удалось получить пользователя из БД: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseBody := map[string]interface{}{"groups": groups, "consists-of": user.ConsistsOf}

	responseBodyJSON, err := json.Marshal(responseBody)
	if err != nil {
		logger.Log.Error("Не удалось выполнить маршалинг JSON: ", err)
		http.Error(w, "Не удалось выполнить маршалинг JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Запишите JSON в тело ответа
	_, err = w.Write(responseBodyJSON)
	if err != nil {
		logger.Log.Error("Не удалось записать JSON в тело ответа: ", err)
		// Здесь мы уже начали отправлять ответ, поэтому можем только залогировать ошибку
		return
	}
}

func (wa *WebApp) HandleGetGroupSchedule(w http.ResponseWriter, r *http.Request) {

	initData, ok := r.Context().Value(InitDataCtx).(*initdata.InitData)
	if !ok || initData == nil {
		logger.Log.Error("Не удалось извлечь initData из контекста")
		http.Error(w, "Не удалось извлечь initData из контекста", http.StatusInternalServerError)
		return
	}

	var jsonData map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil && err != io.EOF {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		logger.Log.Error("Error parsing JSON: " + err.Error())
		return
	}
	defer r.Body.Close()

	initDataInterface, ok := jsonData["group-id"]
	if !ok {
		http.Error(w, "The map does not contain group-id", http.StatusBadRequest)
		logger.Log.Error("The map does not contain group-id")
		return
	}
	groupID, ok := initDataInterface.(uint)
	if !ok {
		http.Error(w, "Can not convert groupID", http.StatusBadRequest)
		logger.Log.Error("Can not convert groupID: " + err.Error())
		return
	}

	group, err := models.DataBase.GetGroup(groupID)
	if err != nil {
		http.Error(w, "Не удалось извлечь группу из БД", http.StatusBadRequest)
		logger.Log.Error("Не удалось извлечь группу из БД: " + err.Error())
		return
	}

	responseBodyJSON, err := json.Marshal(group)
	if err != nil {
		logger.Log.Error("Не удалось выполнить маршалинг JSON: ", err)
		http.Error(w, "Не удалось выполнить маршалинг JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Запишите JSON в тело ответа
	_, err = w.Write(responseBodyJSON)
	if err != nil {
		logger.Log.Error("Не удалось записать JSON в тело ответа: ", err)
		// Здесь мы уже начали отправлять ответ, поэтому можем только залогировать ошибку
		return
	}
}

type SetGroupScheduleRequest struct {
	Type    string `json:"type"`
	GroupId int    `json:"group_id"`

	HolderId int `json:"holder_id"`

	Name  string `json:"name"`
	Title string `json:"title"`

	IsAlternatingGroup bool   `json:"is_alternating_group"`
	StartDate          string `json:"start_date"`

	EvenWeek models.WeekSchedule `json:"even_week"`
	OddWeek  models.WeekSchedule `json:"odd_week"`
}

func (wa *WebApp) HandleSetGroupSchedule(w http.ResponseWriter, r *http.Request) {

	initData, ok := r.Context().Value(InitDataCtx).(*initdata.InitData)
	if !ok || initData == nil {
		logger.Log.Error("Не удалось извлечь initData из контекста")
		http.Error(w, "Не удалось извлечь initData из контекста", http.StatusInternalServerError)
		return
	}

	var jsonData SetGroupScheduleRequest

	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil && err != io.EOF {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		logger.Log.Error("Error parsing JSON: " + err.Error())
		return
	}
	defer r.Body.Close()

	if jsonData.Type != "save" && jsonData.Type != "create" {
		logger.Log.Error("type должен быть только save или create: ", jsonData.Type)
		http.Error(w, "type должен быть только save или create: "+jsonData.Type, http.StatusInternalServerError)
		return
	}

	layout := "2006-01-02"

	t, err := time.Parse(layout, jsonData.StartDate)
	if err != nil {
		logger.Log.Error("Не удалось конвертировать дату: ", err)
		http.Error(w, "Не удалось конвертировать дату: "+err.Error(), http.StatusInternalServerError)
		return
	}

	evenWeekStr, err := json.Marshal(jsonData.EvenWeek)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusBadRequest)
		logger.Log.Error("Error marshaling JSON: " + err.Error())
		return
	}

	oddWeekStr, err := json.Marshal(jsonData.OddWeek)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusBadRequest)
		logger.Log.Error("Error marshaling JSON: " + err.Error())
		return
	}

	group := models.Group{
		HolderID: uint(jsonData.HolderId),
		Name:     jsonData.Name,
		Title:    jsonData.Title,

		StartDate:          &t,
		IsAlternatingGroup: jsonData.IsAlternatingGroup,
		EvenWeek:           string(evenWeekStr),
		OddWeek:            string(oddWeekStr),
	}

	if jsonData.Type == "save" {
		group.ID = uint(jsonData.GroupId)
		err = models.DataBase.UpdateGroup(&group)
		if err != nil {
			http.Error(w, "Не удалось обновить группу в БД: "+err.Error(), http.StatusBadRequest)
			logger.Log.Error("Не удалось обновить группу в БД: " + err.Error())
			return
		}
	}

	if jsonData.Type == "create" {
		_, err = models.DataBase.CreateGroup(&group)
		if err != nil {
			http.Error(w, "Не удалось обновить группу в БД: "+err.Error(), http.StatusBadRequest)
			logger.Log.Error("Не удалось обновить группу в БД: " + err.Error())
			return
		}
	}

	responseBody := map[string]interface{}{"status": true, "message": "Группа успешно создана"}

	responseBodyJSON, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "Error marshaling JSON: "+err.Error(), http.StatusBadRequest)
		logger.Log.Error("Error marshaling JSON: " + err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Запишите JSON в тело ответа
	_, err = w.Write(responseBodyJSON)
	if err != nil {
		logger.Log.Error("Не удалось записать JSON в тело ответа: ", err)
		// Здесь мы уже начали отправлять ответ, поэтому можем только залогировать ошибку
		return
	}
}

func (wa *WebApp) HandleGroupJoin(w http.ResponseWriter, r *http.Request) {

	initData, ok := r.Context().Value(InitDataCtx).(*initdata.InitData)
	if !ok || initData == nil {
		logger.Log.Error("Не удалось извлечь initData из контекста")
		http.Error(w, "Не удалось извлечь initData из контекста", http.StatusInternalServerError)
		return
	}

	var jsonData map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil && err != io.EOF {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		logger.Log.Error("Error parsing JSON: " + err.Error())
		return
	}
	defer r.Body.Close()

	joinToInterface, ok := jsonData["join-to"]
	if !ok {
		http.Error(w, "The map does not contain join-to", http.StatusBadRequest)
		logger.Log.Error("The map does not contain join-to")
		return
	}

	joinToId, ok := joinToInterface.(uint)
	if !ok {
		http.Error(w, "Can not convert joinToId to uint", http.StatusBadRequest)
		logger.Log.Error("Can not convert joinToId to uint: " + err.Error())
		return
	}

	err = models.DataBase.JoinGroup(uint(initData.User.ID), joinToId)
	if err != nil {
		http.Error(w, "Невозможно присоединить пользователя к группе: "+err.Error(), http.StatusBadRequest)
		logger.Log.Error("Невозможно присоединить пользователя к группе: " + err.Error())
		return
	}

	responseBody := map[string]interface{}{"status": true, "message": "Группа успешно создана"}

	responseBodyJSON, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "Error marshaling JSON: "+err.Error(), http.StatusBadRequest)
		logger.Log.Error("Error marshaling JSON: " + err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Запишите JSON в тело ответа
	_, err = w.Write(responseBodyJSON)
	if err != nil {
		logger.Log.Error("Не удалось записать JSON в тело ответа: ", err)
		// Здесь мы уже начали отправлять ответ, поэтому можем только залогировать ошибку
		return
	}
}

func (wa *WebApp) HandleExitGroup(w http.ResponseWriter, r *http.Request) {

	initData, ok := r.Context().Value(InitDataCtx).(*initdata.InitData)
	if !ok || initData == nil {
		logger.Log.Error("Не удалось извлечь initData из контекста")
		http.Error(w, "Не удалось извлечь initData из контекста", http.StatusInternalServerError)
		return
	}

	err := models.DataBase.ExitGroup(uint(initData.User.ID))
	if err != nil {
		http.Error(w, "Невозможно присоединить пользователя к группе: "+err.Error(), http.StatusBadRequest)
		logger.Log.Error("Невозможно присоединить пользователя к группе: " + err.Error())
		return
	}

	responseBody := map[string]interface{}{"status": true, "message": "Группа успешно создана"}

	responseBodyJSON, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "Error marshaling JSON: "+err.Error(), http.StatusBadRequest)
		logger.Log.Error("Error marshaling JSON: " + err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Запишите JSON в тело ответа
	_, err = w.Write(responseBodyJSON)
	if err != nil {
		logger.Log.Error("Не удалось записать JSON в тело ответа: ", err)
		// Здесь мы уже начали отправлять ответ, поэтому можем только залогировать ошибку
		return
	}
}
