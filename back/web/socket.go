package web

import (
	"app/db"
	"app/log"
	"app/model"
	"encoding/json"

	"github.com/gorilla/websocket"
)

const (
	ActionSearchData  = "ActionSearchData"  // Получение данных для поиска (предметы, уровни и т.д.)
	ActionTeacherData = "ActionTeacherData" // Получение списка преподавателей по фильтру
)

// HandleTeacherConnection обрабатывает WebSocket соединение для поиска преподавателей
func HandleTeacherConnection(conn *websocket.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.App.Error("Panic in WebSocket handler:", r)
		}
		conn.Close()
	}()

	for {
		// Читаем сырое сообщение
		_, message, err := conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.App.Error("WebSocket unexpected close:", err)
			break
		}
		if err != nil {
			log.App.Error("Error reading WebSocket message:", err)
			break
		}

		// Сначала парсим как базовый запрос
		var request model.WebSocketRequest
		if err := json.Unmarshal(message, &request); err != nil {
			log.App.Error("Error parsing WebSocket request:", err)
			continue
		}

		log.App.Info(conn.RemoteAddr(), "Received WebSocket request:", request.Action)

		switch request.Action {
		case ActionSearchData:
			log.App.Info(conn.RemoteAddr(), " Getting search data")
			handleSearchDataRequest(conn)
		case ActionTeacherData:
			log.App.Info(conn.RemoteAddr(), " Getting teacher data")
			// Используем те же сырые данные для парсинга TeacherRequest
			var teacherReq model.TeacherSocketRequest
			if err := json.Unmarshal(message, &teacherReq); err != nil {
				log.App.Error("Error parsing teacher request:", err)
				sendError(conn, "Error parsing teacher request", err)
				continue
			}
			handleTeacherDataRequest(conn, teacherReq.Filter)
		default:
			sendError(conn, "Unknown action: "+request.Action, nil)
		}
	}
}

// handleSearchDataRequest обрабатывает запрос на получение данных для поиска
func handleSearchDataRequest(conn *websocket.Conn) {
	log.App.Info(conn.RemoteAddr(), " Getting search data")
	data, err := db.GetSearchData()
	if err != nil {
		sendError(conn, "Error getting search data", err)
		return
	}

	if err := conn.WriteJSON(data); err != nil {
		log.App.Error("Error writing search data:", err)
	}
	log.App.Info("Search data sent successfully")
}

// handleTeacherDataRequest обрабатывает запрос на получение списка преподавателей
func handleTeacherDataRequest(conn *websocket.Conn, filter model.TeacherFilter) {
	log.App.Info("Getting teachers with filter:", filter)

	teachers, total, err := db.GetTeachersByFilter(filter)
	if err != nil {
		log.App.Error("Failed to get teachers:", err)
		sendError(conn, "Error getting teachers", err)
		return
	}

	log.App.Info("Found teachers:", len(teachers), "total:", total)

	response := model.TeachersResponse{
		WebSocketResponse: model.WebSocketResponse{
			Status:  "success",
			Message: "Teachers fetched successfully",
		},
		Teachers: teachers,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}

	if err := conn.WriteJSON(response); err != nil {
		log.App.Error("Failed to write teachers data:", err)
		return
	}
	log.App.Info("Teachers data sent successfully")
}

// sendError отправляет сообщение об ошибке клиенту
func sendError(conn *websocket.Conn, msg string, err error) {
	errMsg := msg
	if err != nil {
		errMsg += ": " + err.Error()
	}

	log.App.Error(errMsg)

	if err := conn.WriteJSON(model.WebSocketResponse{
		Status:  "error",
		Message: errMsg,
	}); err != nil {
		log.App.Error("Error sending error message:", err)
	}
}
