package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	// Step 1: Authenticate User
	fmt.Println("Authenticating user...")
	authResponse, err := makeRequest("POST", "http://localhost:8080/validate-session", map[string]string{
		"token": "dummy-token",
	})
	if err != nil || !authResponse["valid"].(bool) {
		fmt.Println("User authentication failed:", err)
		return
	}
	fmt.Println("User authenticated.")

	// Step 2: Fetch Products
	fmt.Println("Fetching products...")
	products, err := makeRequest("GET", "http://localhost:8082/products", nil)
	if err != nil {
		fmt.Println("Error fetching products:", err)
		return
	}
	fmt.Printf("Products available: %v\n", products)

	// Step 3: Add Product to Cart
	fmt.Println("Adding product to cart...")
	cartItem := map[string]interface{}{
		"product_id": "1",
		"quantity":   1,
	}
	_, err = makeRequest("POST", "http://localhost:8083/cart?user_id=123", cartItem)
	if err != nil {
		fmt.Println("Error adding item to cart:", err)
		return
	}
	fmt.Println("Item added to cart.")

	// Step 4: Checkout
	fmt.Println("Placing order...")
	order := map[string]interface{}{
		"user_id": "123",
		"items": []map[string]interface{}{
			{
				"product_id": "1",
				"quantity":   1,
			},
		},
	}
	orderResponse, err := makeRequest("POST", "http://localhost:8083/checkout", order)
	if err != nil {
		fmt.Println("Error placing order:", err)
		return
	}
	fmt.Printf("Order placed successfully: %v\n", orderResponse)

	// Step 5: Email Notification
	fmt.Println("Triggering email notification...")
	email := map[string]string{
		"to":      "customer@example.com",
		"subject": "Order Confirmation",
		"body":    "Thank you for your order!",
	}
	_, err = makeRequest("POST", "http://localhost:8081/send-email", email)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}
	fmt.Println("Email notification sent.")
}

func makeRequest(method, url string, payload interface{}) (map[string]interface{}, error) {
	var req *http.Request
	var err error

	if payload != nil {
		data, _ := json.Marshal(payload)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(data))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}
