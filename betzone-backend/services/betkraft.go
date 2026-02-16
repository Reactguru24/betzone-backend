package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/betzone/backend/config"
	"github.com/betzone/backend/models"
	"github.com/betzone/backend/utils"
)

type BetkraftService struct {
	baseURL string
	apiKey  string
	appKey  string
	client  *http.Client
}

func NewBetkraftService(cfg *config.Config) *BetkraftService {
	return &BetkraftService{
		baseURL: cfg.BetkraftBaseURL,
		apiKey:  cfg.BetkraftAPIKey,
		appKey:  cfg.BetkraftAppKey,
		client:  &http.Client{},
	}
}

func (bs *BetkraftService) GetGames(page, perPage int, status int) (*models.BetkraftGameResponse, error) {
	endpoint := fmt.Sprintf("%s/v1/games", bs.baseURL)

	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	queryParams := reqURL.Query()
	if page > 0 {
		queryParams.Add("page", fmt.Sprintf("%d", page))
	}
	if perPage > 0 {
		queryParams.Add("per_page", fmt.Sprintf("%d", perPage))
	}
	if status >= 0 {
		queryParams.Add("status", fmt.Sprintf("%d", status))
	}
	reqURL.RawQuery = queryParams.Encode()

	body, err := bs.doRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	var response models.BetkraftGameResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse games response: %v", err)
	}

	fmt.Printf("Successfully fetched %d games from page %d\n", len(response.Data.Data), response.Data.Page)
	return &response, nil
}

func (bs *BetkraftService) GetGameByID(gameID string) (*models.Game, error) {
	endpoint := fmt.Sprintf("%s/v1/games/%s", bs.baseURL, gameID)

	body, err := bs.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var game models.Game
	if err := json.Unmarshal(body, &game); err != nil {
		return nil, err
	}

	return &game, nil
}

func (bs *BetkraftService) LaunchGame(launchReq *models.LaunchGameRequest) (*models.LaunchGameResponse, error) {
	endpoint := fmt.Sprintf("%s/v1/launch", bs.baseURL)

	payload := map[string]interface{}{
		"player_id":    launchReq.PlayerID,
		"player_name":  launchReq.PlayerName,
		"player_token": launchReq.PlayerToken,
		"game_uuid":    launchReq.GameUUID,
		"currency":     launchReq.Currency,
		"balance":      launchReq.Balance,
		"demo":         launchReq.Demo,
	}

	body, err := bs.doRequest("POST", endpoint, payload)
	if err != nil {
		return nil, err
	}

	var response models.LaunchGameResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse launch response: %v", err)
	}

	fmt.Printf("Successfully launched game with UUID: %s\n", launchReq.GameUUID)
	return &response, nil
}

func (bs *BetkraftService) doRequest(method, endpoint string, payload map[string]interface{}) ([]byte, error) {
	var req *http.Request
	var err error

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	timestamp := utils.GetTimestamp()

	var signatureKey string

	if method == "GET" {
		queryParams := u.Query()
		if len(queryParams) > 0 {
			requestMap := make(map[string]interface{})
			for k, v := range queryParams {
				if len(v) > 0 {
					requestMap[k] = v[0]
				}
			}
			signatureKey = utils.HashCreate(requestMap, bs.appKey)
		} else {
			signatureKey = utils.HashCreate(map[string]interface{}{}, bs.appKey)
		}

		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	} else {
		if payload != nil {
			signatureKey = utils.HashCreate(payload, bs.appKey)
		} else {
			signatureKey = utils.HashCreate(map[string]interface{}{}, bs.appKey)
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, u.String(), bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", bs.apiKey)
	req.Header.Set("x-timestamp", timestamp)
	req.Header.Set("x-signature-key", signatureKey)

	fmt.Printf("Request Details:\n")
	fmt.Printf("  Method: %s\n", method)
	fmt.Printf("  Endpoint: %s\n", endpoint)
	fmt.Printf("  API Key (x-api-key): %s\n", bs.apiKey)
	fmt.Printf("  Timestamp: %s\n", timestamp)
	fmt.Printf("  Signature Key: %s\n", signatureKey)
	if len(u.Query()) > 0 {
		fmt.Printf("  Query Params: %s\n", u.Query().Encode())
	}
	fmt.Printf("\n")

	resp, err := bs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n\n", string(body))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %s (status code: %d)", string(body), resp.StatusCode)
	}

	return body, nil
}
