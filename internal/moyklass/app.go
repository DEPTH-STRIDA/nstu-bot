package moyklass

import (
	"app/internal/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"
	"time"
)

const (
	APIUrl = "https://api.moyklass.com"
)

type MoyClassKit struct {
	mux        *http.ServeMux
	updates    chan *Update
	httpClient *http.Client

	apiKey     string
	tokenMutex sync.RWMutex
	token      *TokenInfo
}

type TokenInfo struct {
	AccessToken string
	ExpiresAt   time.Time
}

// NewMoyClassKit создает новый экземпляр MoyClassKit
func NewMoyClassKit(mux *http.ServeMux, webhoukUrl string, apiKey string) (*MoyClassKit, error) {
	webhoukUrl, err := validateRoutePath(webhoukUrl)
	if err != nil {
		return nil, err
	}

	app := &MoyClassKit{
		mux:        mux,
		updates:    make(chan *Update),
		httpClient: &http.Client{},
		apiKey:     apiKey,
	}

	mux.HandleFunc(webhoukUrl, app.HandlerWebHouk)

	if err := app.updateToken(); err != nil {
		return nil, fmt.Errorf("failed to get initial token: %v", err)
	}

	go app.tokenRefresher()

	return app, nil
}

func (app *MoyClassKit) updateToken() error {
	resp, err := app.GetToken(GetTokenConfig{ApiKey: app.apiKey})
	if err != nil {
		return err
	}

	app.tokenMutex.Lock()
	app.token = &TokenInfo{
		AccessToken: resp.AccessToken,
		ExpiresAt:   resp.ExpiresAt,
	}
	app.tokenMutex.Unlock()

	return nil
}

func (app *MoyClassKit) tokenRefresher() {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		if err := app.updateToken(); err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to refresh token: %v", err))
		} else {
			logger.Log.Info("Token refreshed successfully")
		}
	}
}

func (app *MoyClassKit) getAccessToken() string {
	app.tokenMutex.RLock()
	defer app.tokenMutex.RUnlock()
	if app.token == nil {
		return ""
	}
	return app.token.AccessToken
}

func (app *MoyClassKit) doRequest(method, endpoint string, body interface{}, needsAuth bool) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, APIUrl+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if needsAuth {
		token := app.getAccessToken()
		if token == "" {
			return nil, fmt.Errorf("no access token available")
		}
		req.Header.Set("x-access-token", token)
	}

	resp, err := app.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp BadResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %v", err)
		}
		return nil, fmt.Errorf("request failed: %s - %s", errResp.Code, errResp.Message)
	}

	return respBody, nil
}

func (app *MoyClassKit) GetToken(config GetTokenConfig) (*GetTokenResponse, error) {
	respBody, err := app.doRequest(http.MethodPost, GetToken, config, false)
	if err != nil {
		return nil, err
	}

	var response GetTokenResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}

func (app *MoyClassKit) RefreshToken(config RefreshTokenConfig) (*RefreshTokenResponse, error) {
	respBody, err := app.doRequest(http.MethodPost, RefreshToken, nil, true)
	if err != nil {
		return nil, err
	}

	var response RefreshTokenResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}

func (app *MoyClassKit) RevokeToken(config RevokeTokenConfig) (*RevokeTokenResponse, error) {
	respBody, err := app.doRequest(http.MethodPost, RevokeToken, nil, true)
	if err != nil {
		return nil, err
	}

	var response RevokeTokenResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}

func (app *MoyClassKit) GetManagers() ([]GoodManagersResponse, error) {
	respBody, err := app.doRequest(http.MethodGet, Managers, nil, true)
	if err != nil {
		return nil, err
	}

	var response []GoodManagersResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return response, nil
}

func (app *MoyClassKit) GetManager(config ManagerConfig) (*ManagerResponse, error) {
	endpoint := fmt.Sprintf("%s%d", Manager, config.ManagerID)
	respBody, err := app.doRequest(http.MethodGet, endpoint, nil, true)
	if err != nil {
		return nil, err
	}

	var response ManagerResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}

func (app *MoyClassKit) HandlerWebHouk(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to read request body: %v", err))
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Log.Info(string(body))

	var tempUpdate Update
	err = json.Unmarshal(body, &tempUpdate)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to unmarshal request body: %v", err))
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	eventObject, exists := EventTypes[tempUpdate.Event]
	if !exists {
		logger.Log.Error(fmt.Sprintf("Unknown event type: %s", tempUpdate.Event))
		http.Error(w, "Unknown event type", http.StatusBadRequest)
		return
	}

	objType := reflect.TypeOf(eventObject)
	obj := reflect.New(objType).Interface()

	if err := json.Unmarshal(tempUpdate.Object, obj); err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to unmarshal object: %v", err))
		http.Error(w, "Failed to unmarshal object", http.StatusBadRequest)
		return
	}

	parsedUpdate := &Update{
		BaseUpdate: tempUpdate.BaseUpdate,
	}
	parsedUpdate.Object = tempUpdate.Object

	go func() {
		app.updates <- parsedUpdate
	}()

	w.WriteHeader(http.StatusOK)
}

func (app *MoyClassKit) GetUpdatesChan() chan *Update {
	return app.updates
}
