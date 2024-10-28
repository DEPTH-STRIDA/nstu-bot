package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	authBaseURL = "https://login.moyklass.com"
	apiBaseURL  = "https://api.moyklass.com"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	CsrfToken string `json:"csrfToken"`
	ExpiresAt string `json:"expiresAt"`
	Level     string `json:"level"`
}

func main() {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	// Шаг 1: Прямая авторизация через POST /login
	authData := AuthRequest{
		Login:    "Awesome.gail@yandex.ru",
		Password: "2003ivafoV!",
	}

	jsonData, _ := json.Marshal(authData)

	req, _ := http.NewRequest("POST", authBaseURL+"/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://app.moyklass.com")

	fmt.Println("Отправка запроса на авторизацию...")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при авторизации: %v\n", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Printf("Получен ответ с кодом: %d\n", resp.StatusCode)
	fmt.Printf("Тело ответа: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка авторизации. Код: %d, Ответ: %s\n", resp.StatusCode, string(body))
		return
	}

	// Сохраняем cookies и CSRF токен
	cookies := resp.Cookies()
	var loginResp LoginResponse
	json.Unmarshal(body, &loginResp)
	csrfToken := loginResp.CsrfToken

	fmt.Printf("Успешная авторизация! CSRF токен: %s\n", csrfToken)
	fmt.Println("Полученные cookies:")
	for _, cookie := range cookies {
		fmt.Printf("- %s: %s\n", cookie.Name, cookie.Value)
	}

	// Шаг 2: Запрашиваем данные аккаунта
	req, _ = http.NewRequest("GET", apiBaseURL+"/v1/user/baseInfo", nil)
	req.Header.Set("Referer", "https://app.moyklass.com")
	req.Header.Set("x-csrf-token", csrfToken)

	// Добавляем все полученные cookies
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	fmt.Println("Отправка запроса на получение данных аккаунта...")
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при получении данных аккаунта: %v\n", err)
		return
	}
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Printf("Получен ответ с кодом: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка получения данных аккаунта. Код: %d, Ответ: %s\n", resp.StatusCode, string(body))
		return
	}

	fmt.Printf("\nДанные аккаунта получены успешно:\n%s\n", string(body))
}
