# Sprout Go SDK

A Go client library for integrating with the Sprout API. This SDK provides simple methods to create payment orders and securely verify incoming webhooks using both standard `net/http` and the `Gin` framework.

## Installation

To install the package, run the following command in your terminal:

```bash
go get github.com/WhoIsWasay/sprout




Quick Start
1. Initialize the Client
To interact with the Sprout API, you need to initialize a new client using your Secret Key.

package main

import (
	"fmt"
	"[github.com/WhoIsWasay/sprout](https://github.com/WhoIsWasay/sprout)"
)

func main() {
	// Initialize a new Sprout client
	client := sprout.NewClient("your_secret_key_here")
}
