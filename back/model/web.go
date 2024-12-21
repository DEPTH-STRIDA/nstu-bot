package model

type WebConfig struct {
	APPIP   string `envconfig:"APP_IP" default:"localhost"` // IP адрес приложения
	APPPORT string `envconfig:"APP_PORT" default:"8080"`    // Порт приложения

	APPURL            string `envconfig:"APP_URL" default:"http://localhost:8080"` // URL приложения
	APPJWTSecret      string `envconfig:"APP_JWT_SECRET" required:"true"`          // Секретный ключ для JWT
	NumberRepetitions int    `envconfig:"APP_NUMBER_OF_REPETITIONS" default:"15"`  // Количество повторов запроса на замену-перенос
	RepeatPause       int    `envconfig:"APP_REPEAT_PAUSE" default:"15"`           // Пауза между повторами запроса на замену-перенос

	JWTTokenPassword string `envconfig:"JWT_TOKEN_PASSWORD" required:"true"` // Пароль для шифрования JWT токена
}

// Response представляет стандартный ответ
type Response struct {
	Status  string `json:"status"`  // Статус ответа
	Message string `json:"message"` // Описание ответа
}

// RoleResponse представляет ответ с JWT токеном
type RoleResponse struct {
	Role string `json:"role"`
	Name string `json:"name"`
	Response
}

type CodeRequest struct {
	Code string `json:"code"`
}

type EmailRequest struct {
	Email string `json:"email"`
}

type PasswordRequest struct {
	Password string `json:"password"`
}

type EmailPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TeacherIDRequest struct {
	TeacherID uint `json:"teacher_id"`
}

type TeacherIDResponse struct {
	Response
	Teacher Teacher `json:"teacher"`
}

type UpdateUserRequest struct {
	Name    string  `json:"name"`
	Student Student `json:"student"`
	Teacher Teacher `json:"teacher"`
	Admin   Admin   `json:"admin"`
}

type UserUpdateRequest struct {
	Admin   Admin   `json:"admin"`
	Teacher Teacher `json:"teacher"`
	Student Student `json:"student"`
}

type AdminResponse struct {
	Response
	Admin Admin `json:"admin"`
}

type TeacherResponse struct {
	Response
	Teacher Teacher `json:"teacher"`
}

type StudentResponse struct {
	Response
	Student Student `json:"student"`
}
