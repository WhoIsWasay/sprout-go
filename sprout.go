package sprout

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
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


//Creates Order Request and send it to /ordercreation and send back reply to the caller
func (cli *Client) CreateOrder(or OrderRequest) (*OrderResponse, error) {

	// 1. Prepares your JSON body
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

type WebhookPayload struct {
	Event     string      `json:"event"`
	Timestamp string      `json:"timestamp"`
	Data      WebhookData `json:"data"`
}

type WebhookData struct {
	Order_id        string `json:"order_id"`
	Merchant_id     string `json:"merchant_id"` 
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	Customer_email  string `json:"customer_email"`
	Idempotency_key string `json:"idempotency_key"`
	Mode            string `json:"mode"`
}


func WebhookReader(secret string, r *http.Request) (*WebhookPayload, error) {
    body, _ := io.ReadAll(r.Body)

    // Verify signature
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(body)
    expected := hex.EncodeToString(mac.Sum(nil))
    received := r.Header.Get("X-Sprout-Signature")

    if !hmac.Equal([]byte(expected), []byte(received)) {
        return nil, fmt.Errorf("invalid signature")
    }

    var payload WebhookPayload
    if err := json.Unmarshal(body, &payload); err != nil {
        return nil, err
    }

    return &payload, nil
}


func WebhookReaderGin(secret string, c *gin.Context) (*WebhookPayload, error) {
    body, _ := io.ReadAll(c.Request.Body)

    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(body)
    expected := hex.EncodeToString(mac.Sum(nil))
    received := c.GetHeader("X-Sprout-Signature")

    if !hmac.Equal([]byte(expected), []byte(received)) {
        return nil, fmt.Errorf("invalid signature")
    }

    var payload WebhookPayload
    if err := json.Unmarshal(body, &payload); err != nil {
        return nil, err
    }

    return &payload, nil
}