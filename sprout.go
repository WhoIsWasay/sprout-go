package sprout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	SecretKey  string
	BaseURL    string // defaults to https://pay.sprout.pk
	HTTPClient *http.Client
}

type CustomerData struct {
	Email string `json:"customer-email"`
}

type OrderResponse struct {
	Status    bool
	SessionID string
}

func NewClient(secretKey string) *Client {
	return &Client{
		SecretKey:  secretKey,
		BaseURL:    "http://localhost:8080/ordercreation",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}


type OrderCreation struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type OrderRequest struct {
	OrderID  string       `json:"orderId"`
	Currency string       `json:"currency"`
	Amount   string       `json:"amount"`
	Customer CustomerData `json:"customer-data"`
}

func (cli *Client) CreateOrder(or OrderRequest) (*OrderResponse, error) {

	// 1. Prepare your JSON body
	jsonBody, err := json.Marshal(or)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(jsonBody)
	// 2. Create a new request object
	req, err := http.NewRequest("POST", cli.BaseURL, bodyReader)
	if err != nil {
		return nil, err
	}

	// 3. ADD YOUR HEADERS HERE
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Secret-Key", cli.SecretKey)

	// 4. Send the request using an HTTP Client
	resp, err := cli.HTTPClient.Do(req) // .Do() executes the request object
	if err != nil {
		return nil, err //External site unreachable

	}
	defer resp.Body.Close()


	// Read the response from the server
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sprout error: %s", string(respBody))
	}

	var result map[string]any

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	status, _ := result["status"].(bool)
	sessionID, _ := result["session-id"].(string)

	orderResponse := OrderResponse{
		Status:    status,
		SessionID: sessionID,
	}

	// RETURN the result back to the caller
	return &orderResponse, nil

}
