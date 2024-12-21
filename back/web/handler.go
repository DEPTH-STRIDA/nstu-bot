package web

import (
	"app/auth"
	"app/db"
	"app/log"
	"app/model"
	"app/utils"
	"encoding/json"
	"net/http"
)

const secureCookie = false

// HandleConnections обрабатывает входящие соединения WebSocket
// @Summary Установка WebSocket соединения
// @Description Устанавливает WebSocket соединение для обмена данными в реальном времени. "Try it out" не поддерживается для WebSocket.
// @Tags websocket
// @Success 101 {object} model.Response "WebSocket соединение успешно установлено"
// @Failure 400 {object} model.Response "Ошибка при установке соединения"
// @Router /ws [get]
func (app *WebApp) HandleConnections(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	log.App.Info("Новое подключение от: ", clientIP)

	ws, err := app.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.App.Error("Ошибка при установке соединения: ", err)
		response := model.Response{Status: "error", Message: "Ошибка при установке соединения"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer ws.Close()

	response := model.Response{Status: "success", Message: "WebSocket соединение успешно установлено"}
	w.WriteHeader(http.StatusSwitchingProtocols)
	json.NewEncoder(w).Encode(response)

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.App.Error("Error reading message:", err)
			break
		}
		log.App.Error("Received message of type", messageType, " : ", message)

		if err := ws.WriteMessage(messageType, message); err != nil {
			log.App.Error("Error writing message:", err)
			break
		}
	}
}

// HandleRegistrationStarted обрабатывает начало регистрации
// @Summary Начало регистрации
// @Description Обрабатывает запрос на начало регистрации пользователя
// @Tags registration
// @Accept json
// @Produce json
// @Param request body model.EmailPasswordRequest true "Запрос на начало регистрации"
// @Success 200 {object} model.Response "Регистрация начата"
// @Failure 400 {object} model.Response "Ошибка в запросе"
// @Router /registration/start [post]
func (app *WebApp) HandleRegistrationStarted(w http.ResponseWriter, r *http.Request) {
	var req model.EmailPasswordRequest

	// Декодируем JSON из тела запроса в структуру
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.App.Error(r.RemoteAddr, " failed to decode registration request: ", err)
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Не удалось распарсить запрос: " + err.Error()}), http.StatusBadRequest)
		return
	}

	response, token, err := auth.RegisterStudent(req.Email, req.Password)
	if err != nil {
		log.App.Error(r.RemoteAddr, " failed to register student: ", err)
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: err.Error()}), http.StatusBadRequest)
		return
	}

	// Устанавливаем токен в httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "registrationToken",
		Value:    token,
		HttpOnly: true,
		Secure:   secureCookie,
		Path:     "/",
	})

	log.App.Info("Received registration request for email: ", req.Email)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleValidateEmail проверяет доступность email для регистрации
// @Summary Проверка доступности email для регистрации
// @Description Обрабатывает запрос на проверку доступности email для регистрации
// @Tags registration
// @Accept json
// @Produce json
// @Param request body model.EmailRequest true "Запрос на проверку доступности email для регистрации"
// @Success 200 {object} model.Response "Email доступен для регистрации"
// @Failure 400 {object} model.Response "Email недоступен для регистрации"
// @Router /registration/validate-email [post]
func (app *WebApp) HandleValidateEmail(w http.ResponseWriter, r *http.Request) {
	var req model.EmailRequest

	// Декодируем JSON из тела запроса в структуру
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Не удалось распарсить запрос: " + err.Error()}), http.StatusBadRequest)
		return
	}

	err = auth.ValidateEmail(req.Email)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: err.Error()}), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{Status: "success", Message: "Email доступен для регистрации"})
}

// HandleRegistrationConfirmation обрабатывает подтверждение регистрации
// @Summary Подтверждение регистрации
// @Description Обрабатывает запрос на подтверждение регистрации пользователя
// @Tags registration
// @Accept json
// @Produce json
// @Param request body model.CodeRequest true "Запрос на подтверждение регистрации"
// @Success 200 {object} model.RoleResponse "Регистрация подтверждена"
// @Failure 400 {object} model.Response "Ошибка в запросе"
// @Router /registration/confirm [post]
func (app *WebApp) HandleRegistrationConfirmation(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("registrationToken")
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Отсутствует токен регистрации"}), http.StatusBadRequest)
		return
	}
	token := cookie.Value

	var req model.CodeRequest

	// Декодируем JSON из тела запроса в структуру
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Не удалось распарсить запрос: " + err.Error()}), http.StatusBadRequest)
		return
	}

	response, token, err := auth.ConfirmRegistration(token, req.Code)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: err.Error()}), http.StatusBadRequest)
		return
	}

	// Устанавливаем токен в httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    token,
		HttpOnly: true,
		Secure:   secureCookie,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleLogin обрабатывает вход пользователя
// @Summary Вход пользователя
// @Description Обрабатывает запрос на вход пользователя. При успешной аутентификации JWT токен будет отправлен в httpOnly cookie и автоматически использоваться в последующих запросах.
// @Tags login
// @Accept json
// @Produce json
// @Param request body model.EmailPasswordRequest true "Запрос на вход"
// @Success 200 {object} model.RoleResponse "Вход выполнен, токен охранен в cookie"
// @Failure 400 {object} model.Response "Ошибка в запросе"
// @Router /auth/login [post]
func (app *WebApp) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req model.EmailPasswordRequest

	// Декодируем JSON из тела запроса в структуру
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.App.Error(r.RemoteAddr, " failed to decode login request: ", err)
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Не удалось распарсить запрос: " + err.Error()}), http.StatusBadRequest)
		return
	}

	response, token, err := auth.Login(req.Email, req.Password)
	if err != nil {
		log.App.Error(r.RemoteAddr, " failed to login: ", err)
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: err.Error()}), http.StatusBadRequest)
		return
	}

	// Устанавливаем токен в httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    token,
		HttpOnly: true,
		Secure:   secureCookie,
		Path:     "/",
	})

	log.App.Info("User logged in successfully: ", req.Email)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleJwtLogin обрабатывает вход пользователя по JWT токену
// @Summary Вход по JWT токену
// @Description Проверяет валидность JWT токена и выполняет вход пользователя
// @Tags login
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "JWT токен в формате: Bearer <token>"
// @Success 200 {object} model.RoleResponse "Вход выполнен успешно"
// @Failure 401 {object} model.Response "Недействительный токен"
// @Router /auth/jwt [post]
func (app *WebApp) HandleJwtLogin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Отсутствует токен авторизации"}), http.StatusBadRequest)
		return
	}
	token := cookie.Value

	response, err := auth.JwtLogin(token)
	if err != nil {
		log.App.Info(r.RemoteAddr, " is not auth")
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Invalid token"}), http.StatusUnauthorized)
		return
	}

	log.App.Info(r.RemoteAddr, " is auth")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleLogout обрабатывает выход пользователя
// @Summary Выход пользователя
// @Description Удаляет токен авторизации из cookie и завершает сессию пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} model.Response "Выход выполнен успешно"
// @Failure 400 {object} model.Response "Ошибка при выходе"
// @Router /auth/logout [post]
func (app *WebApp) Logout(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем токен в httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    "",
		HttpOnly: true,
		Secure:   secureCookie,
		Path:     "/",
		MaxAge:   -1,
	})

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{Status: "success", Message: "Выход выполнен успешно"})
}

// HandleResetPasswordStarted обрабатывает начало сброса пароля
// @Summary Начало сброса пароля
// @Description Обрабатывает запрос на начало сброса пароля
// @Tags password
// @Accept json
// @Produce json
// @Param request body model.EmailPasswordRequest true "Запрос на начало сброса пароля"
// @Success 200 {object} model.Response "Сброс пароля начат"
// @Failure 400 {object} model.Response "Ошибка в запросе"
// @Router /password/reset/start [post]
func (app *WebApp) HandleResetPasswordStarted(w http.ResponseWriter, r *http.Request) {
	var req model.EmailPasswordRequest

	// Декодируем JSON из тела запроса в структуру
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Не удалось распарсить запрос: " + err.Error()}), http.StatusBadRequest)
		return
	}

	token, err := auth.ResetPassword(req.Email)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: err.Error()}), http.StatusBadRequest)
		return
	}

	// Устанавливаем токен в httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "resetPasswordToken",
		Value:    token,
		HttpOnly: true,
		Secure:   secureCookie,
		Path:     "/",
	})

	log.App.Info("Received registration request for email: ", req.Email)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{Status: "success", Message: "Сброс пароля начат. На почту отправлен код подтверждения."})
}

// HandleResetPasswordConfirmation обрабатывает подтверждение сброса пароля
// @Summary Подтверждение сброса пароля
// @Description Обрабатывает запрос на подтверждение сброса пароля
// @Tags password
// @Accept json
// @Produce json
// @Param request body model.CodeRequest true "Запрос на подтверждение сброса пароля"
// @Success 200 {object} model.Response "Сброс пароля подтвержден"
// @Failure 400 {object} model.Response "Ошибка в запросе"
// @Router /password/reset/confirm [post]
func (app *WebApp) HandleResetPasswordConfirmation(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("resetPasswordToken")
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Отсутствует токен восстановления пароля"}), http.StatusBadRequest)
		return
	}
	token := cookie.Value

	var req model.CodeRequest

	// Декодируем JSON из тела запроса в структуру
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Не удалось распарсить запрос: " + err.Error()}), http.StatusBadRequest)
		return
	}

	token, err = auth.ConfirmResetPassword(token, req.Code)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: err.Error()}), http.StatusBadRequest)
		return
	}

	// Устанавливаем токен в httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "newPasswordToken",
		Value:    token,
		HttpOnly: true,
		Secure:   secureCookie,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{Status: "success", Message: "Сброс пароля подтвержден. Введите новый пароль."})
}

// HandleSetNewPassword обрабатывает установку нового пароля
// @Summary Установка нового пароля
// @Description Обрабатывает запрос на установку нового пароля
// @Tags password
// @Accept json
// @Produce json
// @Param request body model.PasswordRequest true "Запрос на установку нового пароля"
// @Success 200 {object} model.RoleResponse "Новый пароль установлен"
// @Failure 400 {object} model.Response "Ошибка в запросе"
// @Router /password/new [post]
func (app *WebApp) HandleSetNewPassword(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("newPasswordToken")
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Отсутствует токен восстановления пароля"}), http.StatusBadRequest)
		return
	}
	token := cookie.Value

	var req model.PasswordRequest

	// Декодируем JSON из тела запроса в структуру
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Не удалось распарсить запрос: " + err.Error()}), http.StatusBadRequest)
		return
	}

	token, err = auth.SetNewPassword(token, req.Password)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: err.Error()}), http.StatusBadRequest)
		return
	}

	// Устанавливаем токен в httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    token,
		HttpOnly: true,
		Secure:   secureCookie,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{Status: "success", Message: "Новый пароль установлен"})

}

// HandleMailAuthentication обрабатывает аутентификацию по почте
// @Summary Аутентификация по почте
// @Description Обрабатывает запрос на аутентификацию по почте
// @Tags mail
// @Accept json
// @Produce json
// @Param request body model.EmailRequest true "Запрос на аутентификацию по почте"
// @Success 200 {object} model.Response "Аутентификация по почте выполнена"
// @Failure 400 {object} model.Response "Ошибка в запросе"
// @Router /mail/auth [post]
func (app *WebApp) HandleMailAuthentication(w http.ResponseWriter, r *http.Request) {
}

// HandleMailAuthConfirmation обрабатывает подтверждение аутентификации по почте
// @Summary Подтверждение аутентификации по почте
// @Description Обрабатывает запрос на подтверждение аутентификации по почте
// @Tags mail
// @Accept json
// @Produce json
// @Param request body model.CodeRequest true "Запрос на подтверждение аутентификации по почте"
// @Success 200 {object} model.RoleResponse "Аутентификация по почте подтверждена"
// @Failure 400 {object} model.Response "Ошибка в запросе"
// @Router /mail/auth/confirm [post]
func (app *WebApp) HandleMailAuthConfirmation(w http.ResponseWriter, r *http.Request) {
}

// HandleUpdateUser обновляет данные пользователя
// @Summary Обновление данных пользователя
// @Description Обрабатывает запрос на обновление данных пользователя
// @Tags mail
// @Accept json
// @Produce json
// @Param request body model.UpdateUserRequest true "Запрос на обновление данных пользователя"
// @Success 200 {object} model.Response "Данные пользователя обновлены"
// @Failure 400 {object} model.Response "Ошибка в запросе"
// @Router /user/update [post]
func (app *WebApp) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Отсутствует токен авторизации"}), http.StatusBadRequest)
		return
	}
	token := cookie.Value

	var req model.UpdateUserRequest

	// Декодируем JSON из тела запроса в структуру
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Не удалось распарсить запрос: " + err.Error()}), http.StatusBadRequest)
		return
	}

	err = auth.UpdateUser(token, req)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Invalid token"}), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{Status: "success", Message: "Данные пользователя обновлены"})
}

// HandleTeacherSearch обрабатывает WebSocket соединение для поиска преподавателей
// @Summary Поиск преподавателей через WebSocket
// @Description Устанавливает WebSocket соединение для динамического поиска преподавателей с фильтрацией и пагинацией
// @Tags websocket
// @Accept json
// @Produce json
// @Param filter body model.TeacherFilter true "Параметры фильтрации преподавателей"
// @Success 101 {object} model.Response "WebSocket соединение успешно установлено"
// @Failure 400 {object} model.Response "Ошибка при установке соединения"
// @Router /ws/teachers [get]
func (app *WebApp) HandleTeacherSearch(w http.ResponseWriter, r *http.Request) {
	// Разрешаем все источники
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.App.Info(r.RemoteAddr, " Attempting to upgrade connection to WebSocket")
	ws, err := app.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.App.Error(r.RemoteAddr, "Failed to upgrade connection:", err)
		return
	}

	log.App.Info(r.RemoteAddr, "WebSocket connection established successfully")
	defer ws.Close()

	HandleTeacherConnection(ws)
}

func (app *WebApp) HandleGetTeacher(w http.ResponseWriter, r *http.Request) {
	var req model.TeacherIDRequest

	// Декодируем JSON из тела запроса в структуру
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Не удалось распар��ить запрос: " + err.Error()}), http.StatusBadRequest)
		return
	}

	// Создаем переменную для хранения преподавателя
	var teacher model.Teacher

	// Получаем преподавателя с предзагрузкой связанных данных
	err = db.App.
		Preload("Subjects").
		Preload("LevelTraining").
		Preload("Experience").
		Preload("Services").
		Where("id = ?", req.TeacherID).
		First(&teacher).Error

	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Преподаватель не найден"}), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.TeacherIDResponse{
		Response: model.Response{
			Status:  "success",
			Message: "Преподаватель найден",
		},
		Teacher: teacher,
	})
}

func (app *WebApp) HandleGetProfileInfo(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{
			Status:  "failure",
			Message: "Отсутствует токен авторизации",
		}), http.StatusBadRequest)
		return
	}
	token := cookie.Value

	admin, teacher, student, err := auth.GetProfileInfo(token)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{
			Status:  "failure",
			Message: "Ошибка при получении данных профиля: " + err.Error(),
		}), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	switch {
	case admin != nil:
		json.NewEncoder(w).Encode(model.AdminResponse{
			Response: model.Response{
				Status:  "success",
				Message: "Данные профиля получены",
			},
			Admin: *admin,
		})
	case teacher != nil:
		json.NewEncoder(w).Encode(model.TeacherResponse{
			Response: model.Response{
				Status:  "success",
				Message: "Данные профиля получены",
			},
			Teacher: *teacher,
		})
	case student != nil:
		json.NewEncoder(w).Encode(model.StudentResponse{
			Response: model.Response{
				Status:  "success",
				Message: "Данные профиля получены",
			},
			Student: *student,
		})
	default:
		http.Error(w, utils.StructToJSONString(model.Response{
			Status:  "failure",
			Message: "Профиль не найден",
		}), http.StatusNotFound)
	}
}

func (app *WebApp) HandleSaveUserInfo(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{
			Status:  "failure",
			Message: "Отсутствует токен авторизации",
		}), http.StatusBadRequest)
		return
	}
	token := cookie.Value

	var req model.UpdateUserRequest

	// Декодируем JSON из тела запроса в структуру
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: "Не удалось распарсить запрос: " + err.Error()}), http.StatusBadRequest)
		return
	}

	err = auth.SaveUserInfo(token, req)
	if err != nil {
		http.Error(w, utils.StructToJSONString(model.Response{Status: "failure", Message: err.Error()}), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{Status: "success", Message: "Данные пользователя обновлены"})
}

func (app *WebApp) HandleGetSearchData(w http.ResponseWriter, r *http.Request) {
	response := model.DataResponse{
		WebSocketResponse: model.WebSocketResponse{
			Status:  "success",
			Message: "Data fetched successfully",
		},
	}

	var subjects []model.Subject
	var levelTrainings []model.LevelTraining
	var levelExperience []model.TutorExperience

	db.App.Model(&subjects).Find(&subjects)
	db.App.Model(&levelTrainings).Find(&levelTrainings)
	db.App.Model(&levelExperience).Find(&levelExperience)

	response.Subjects = subjects
	response.LevelTrainings = levelTrainings
	response.LevelExperience = levelExperience

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
