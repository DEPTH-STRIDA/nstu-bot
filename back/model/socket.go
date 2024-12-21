package model

// TeacherFilter представляет параметры фильтрации преподавателей
type TeacherFilter struct {
	Name          string `json:"name"`
	Subjects      []uint `json:"subjects"`
	LevelTraining []uint `json:"level_training"`
	Experience    uint   `json:"experience"`
	ClassFormat   string `json:"class_format,omitempty"` // "Индивидуальные", "Групповые", "Индивидуальные и групповые"
	PriceFrom     *int   `json:"price_from,omitempty"`
	PriceTo       *int   `json:"price_to,omitempty"`
	Page          int    `json:"page"`
	PageSize      int    `json:"page_size"`
}

// WebSocketRequest представляет запрос через WebSocket
type WebSocketRequest struct {
	Action string `json:"action"`
}

// WebSocketResponse представляет базовый ответ WebSocket
type WebSocketResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type TeacherSocketRequest struct {
	WebSocketRequest
	Filter TeacherFilter `json:"filter"`
}

type TeachersResponse struct {
	WebSocketResponse
	Teachers []Teacher `json:"teachers"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}

// DataResponse представляет ответ с данными для фильтров
type DataResponse struct {
	WebSocketResponse
	Subjects        []Subject         `json:"subjects"`
	LevelTrainings  []LevelTraining   `json:"level_trainings"`
	LevelExperience []TutorExperience `json:"level_experience"`
	ClassFormats    []ClassFormat     `json:"class_formats"`
}
