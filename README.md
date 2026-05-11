# Sprout Go SDK

A Go client library for integrating with the Sprout API. This SDK provides simple methods to create payment orders and securely verify incoming webhooks using both standard `net/http` and the `Gin` framework.

## Installation

To install the package, run the following command in your terminal:

```bash
go get github.com/WhoIsWasay/sprout-go
```

# Quick Start
### 1. Initialize the Client

To interact with the Sprout API, you need to initialize a new client using your <b>Secret Key</b>.

```bash
package main

import (
	"fmt"
	"github.com/WhoIsWasay/sprout-go"
)

func main() {
	// Initialize a new Sprout client
	client := sprout.NewClient("your_secret_key_here")
}

```

### 2. Create an Order
You can create an order by passing an OrderRequest struct to the CreateOrder method. The SDK handles the JSON marshaling, setting the appropriate headers, and parsing the response.

```bash
package main

import (
	"fmt"
	"log"
	"github.com/WhoIsWasay/sprout-go"
)

func main() {
	client := sprout.NewClient("your_secret_key_here")

	// Prepare the order details
	req := sprout.OrderRequest{
		OrderID:  "ORD-12345",
		Currency: "PKR",
		Amount:   "5000.00",
		Customer: sprout.CustomerData{
			Email: "customer@example.com",
		},
	}

	// Send the request
	response, err := client.CreateOrder(req)
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}

	// Check status and retrieve Session ID
	if response.Status {
		fmt.Printf("Order created successfully! Session ID: %s\n", response.SessionID)
	} else {
		fmt.Println("Order creation failed.")
	}
}

```
# Handling Webhooks

When events happen in Sprout (like a successful payment), Sprout will send an HTTP POST request to your webhook endpoint.

This SDK provides built-in methods to <b>automatically verify the HMAC SHA-256  signature </b>to ensure the request actually came from Sprout, and parses the payload for you.

### Option A: Standard `net/http` Webhook

"If you are using Go's standard `net/http` package, use `WebhookReader`:


```bash
package main

import (
	"fmt"
	"net/http"
	"github.com/WhoIsWasay/sprout-go"
)

func handleSproutWebhook(w http.ResponseWriter, r *http.Request) {
	secret := "your_webhook_secret"
	
	// Verifies signature and parses payload
	payload, err := sprout.WebhookReader(secret, r)
	if err != nil {
		http.Error(w, "Invalid signature or payload", http.StatusUnauthorized)
		return
	}

	fmt.Printf("Received event: %s for Order: %s\n", payload.Event, payload.Data.Order_id)
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/webhook", handleSproutWebhook)
	http.ListenAndServe(":8080", nil)
}
```

### Option B: `Gin` Framework Webhook
If you are using the `Gin` Web Framework, use the dedicated `WebhookReaderGin` method:

```bash
package main

import (
	"net/http"
    "github.com/gin-gonic/gin"
    "github.com/WhoIsWasay/sprout-go"
)

func main() {
	router := gin.Default()

	router.POST("/webhook", func(c *gin.Context) {
		secret := "your_webhook_secret"

		// Verifies signature and parses payload directly from the Gin context
		payload, err := sprout.WebhookReaderGin(secret, c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature or payload"})
			return
		}

		// Handle successful payload
		c.JSON(http.StatusOK, gin.H{
			"message": "Webhook received successfully",
			"event":   payload.Event,
			"orderId": payload.Data.Order_id,
		})
	})

	router.Run(":8080")
}
```
# Types Reference
### OrderRequest

`OrderID` (string): Your internal order reference.

`Currency` (string): The currency code (e.g., "PKR").

`Amount` (string): The transaction amount.

`Customer` (`CustomerData`): Contains the customer's email.

### OrderResponse

`Status` (bool): Indicates if the order was created successfully.

`SessionID` (string): The unique session ID returned by Sprout.

### WebhookPayload

`Event` (string): The type of event triggered.

`Timestamp` (string): Time the event was triggered.

`Data` (`WebhookData`): Contains detailed information about the transaction (amount, currency, mode, idempotency key, etc.).
