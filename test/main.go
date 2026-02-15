package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	nginxHost = "nginx"
	httpPort  = "80"
	httpsPort = "443"
	errorPort = "8080"
)

func main() {
	fmt.Println("=== F5 DevOps Assignment - Nginx Test Suite ===\n")

	allPassed := true

	// Test 1: HTTP Server (port 80)
	if !testHTTPServer() {
		allPassed = false
	}

	// Test 2: HTTPS Server (port 443)
	if !testHTTPSServer() {
		allPassed = false
	}

	// Test 3: Error Server (port 8080)
	if !testErrorServer() {
		allPassed = false
	}

	// Test 4: Rate Limiting
	if !testRateLimiting() {
		allPassed = false
	}

	// Summary
	fmt.Println("\n Test Summary ")
	if allPassed {
		fmt.Println("All tests passed!")
		os.Exit(0)
	} else {
		fmt.Println("Some tests failed")
		os.Exit(1)
	}
}

func testHTTPServer() bool {
	fmt.Println("[Test 1] HTTP Server (port 80)")

	url := fmt.Sprintf("http://%s:%s/", nginxHost, httpPort)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("  FAIL: Could not connect to %s: %v\n", url, err)
		return false
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != 200 {
		fmt.Printf("  FAIL: Expected status 200, got %d\n", resp.StatusCode)
		return false
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("  FAIL: Could not read response body: %v\n", err)
		return false
	}

	// Check if response contains HTML
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "<html") && !strings.Contains(bodyStr, "<HTML") {
		fmt.Printf("  FAIL: Response does not appear to be HTML\n")
		return false
	}

	fmt.Printf("  PASS: Received 200 OK with HTML content\n")
	return true
}

func testHTTPSServer() bool {
	fmt.Println("[Test 2] HTTPS Server (port 443)")

	// Create HTTP client that accepts self-signed certificates
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	url := fmt.Sprintf("https://%s:%s/", nginxHost, httpsPort)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("  FAIL: Could not connect to %s: %v\n", url, err)
		return false
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != 200 {
		fmt.Printf("  FAIL: Expected status 200, got %d\n", resp.StatusCode)
		return false
	}

	fmt.Printf("  PASS: HTTPS server responding with 200 OK\n")
	return true
}

func testErrorServer() bool {
	fmt.Println("[Test 3] Error Server (port 8080)")

	url := fmt.Sprintf("http://%s:%s/", nginxHost, errorPort)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("  FAIL: Could not connect to %s: %v\n", url, err)
		return false
	}
	defer resp.Body.Close()

	// Check status code - should be 403
	if resp.StatusCode != 403 {
		fmt.Printf("  FAIL: Expected status 403, got %d\n", resp.StatusCode)
		return false
	}

	fmt.Printf("  PASS: Error server returning 403 Forbidden\n")
	return true
}

func testRateLimiting() bool {
	fmt.Println("[Test 4] Rate Limiting")

	url := fmt.Sprintf("http://%s:%s/", nginxHost, httpPort)
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// Send rapid requests
	rateLimitHit := false
	successCount := 0
	rateLimitCount := 0

	fmt.Printf("  Sending 20 rapid requests...\n")
	for i := 0; i < 20; i++ {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("  FAIL: Request %d failed: %v\n", i+1, err)
			return false
		}

		if resp.StatusCode == 200 {
			successCount++
		} else if resp.StatusCode == 503 {
			rateLimitHit = true
			rateLimitCount++
		}

		resp.Body.Close()

		// Small delay to avoid overwhelming the system
		time.Sleep(10 * time.Millisecond)
	}

	if !rateLimitHit {
		fmt.Printf("  FAIL: Rate limiting not triggered (all requests succeeded)\n")
		fmt.Printf("     Expected at least one 503 response\n")
		return false
	}

	fmt.Printf("  PASS: Rate limiting working (%d successful, %d rate-limited)\n",
		successCount, rateLimitCount)
	return true
}
