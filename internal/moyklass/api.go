package moyklass

import "time"

type BadResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Meta    string `json:"meta,omitempty"`
}

// ApiKey string `json:"apiKey,omitempty"`

// Авторизация. Получение токена для работы с API
type GetTokenConfig struct {
	ApiKey string `json:"apiKey"`
}

type GetTokenResponse struct {
	AccessToken string    `json:"accessToken,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
	Level       string    `json:"level,omitempty"`
	BadResponse
}

// RefreshTokenConfig генерирует новый токен, текущий токен при этом продолжает действовать.
type RefreshTokenConfig struct {
	XAccessToken string `json:"-"`
}

type RefreshTokenResponse struct {
	AccessToken string    `json:"accessToken,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
	Level       string    `json:"level,omitempty"`
	BadResponse
}

// RevokeTokenConfig удаляет существующий токен. Токен передается в заголовке x-access-token
type RevokeTokenConfig struct {
	XAccessToken string `json:"-"`
}
type RevokeTokenResponse struct {
	BadResponse
}

type ManagersConfig struct {
	XAccessToken string `json:"-"`
}

type ManagerStruct struct {
	ID                 int64   `json:"id"`
	Name               string  `json:"name"`
	Phone              string  `json:"phone"`
	Email              string  `json:"email"`
	Filials            []int64 `json:"filials"`
	SalaryFilialID     int64   `json:"salaryFilialId"`
	Roles              []int64 `json:"roles"`
	Enabled            bool    `json:"enabled"`
	AdditionalContacts string  `json:"additionalContacts"`
	IsStaff            bool    `json:"isStaff"`
	IsWork             bool    `json:"isWork"`
	SendNotifies       bool    `json:"sendNotifies"`
	StartDate          string  `json:"startDate"`
	EndDate            string  `json:"endDate"`
	ContractNumber     string  `json:"contractNumber"`
	ContractDate       string  `json:"contractDate"`
	BirthDate          string  `json:"birthDate"`
	PassportData       string  `json:"passportData"`
	Comment            string  `json:"comment"`
	Color              string  `json:"color"`
	RateID             int64   `json:"rateId"`
	IsOwner            bool    `json:"isOwner"`
	AllowFunds         bool    `json:"allowFunds"`
	LastActive         string  `json:"lastActive"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

type GoodManagersResponse struct {
	ManagerStruct
}

type BadManagersResponse struct {
	BadResponse
}

type ManagerConfig struct {
	ManagerID    int
	XAccessToken string `json:"-"`
}

type ManagerResponse struct {
	ManagerStruct
	BadResponse
}
