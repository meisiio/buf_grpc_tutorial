package main

import (
	"errors"
	"fmt"
	"time"

	gobreaker "github.com/sony/gobreaker/v2"
)

func main() {
	// 1. Configure the Circuit Breaker
	var cb *gobreaker.CircuitBreaker[[]byte]
	cb = gobreaker.NewCircuitBreaker[[]byte](gobreaker.Settings{
		Name:        "PaymentAPI",
		MaxRequests: 1,               // How many requests are allowed in Half-Open state
		Interval:    0,               // Cyclic period of the closed state (0 = never clears counts)
		Timeout:     3 * time.Second, // How long to stay OPEN before switching to HALF-OPEN
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip the breaker if 3 requests fail in a row!
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("\n[CIRCUIT BREAKER] State changed from %s to %s!\n\n", from.String(), to.String())
		},
	})

	// 2. Simulate 20 requests
	for i := 1; i <= 20; i++ {
		// We execute our API call INSIDE the circuit breaker!
		_, err := cb.Execute(func() ([]byte, error) {
			return callFlakyAPI(i)
		})

		if err != nil {
			fmt.Printf("Request %d failed: %v\n", i, err)
		}

		time.Sleep(1 * time.Second) // Wait 1 second between requests
	}
}

// 3. The Flaky API Simulator
func callFlakyAPI(attempt int) ([]byte, error) {
	// TODO: Write logic here so that requests 1 through 5 return an error (simulate a server crash).
	// But starting on request 6, it magically starts returning nil (simulate the server recovering).
	// If it succeeds, print: fmt.Printf("Request %d reached the API and succeeded!\n", attempt)
	// If it fails, print: fmt.Printf("Request %d reached the API and FAILED!\n", attempt)

	if attempt <= 5 {
		fmt.Printf("Request %d reached the API and FAILED!\n", attempt)
		return nil, errors.New("not implemented")
	} else if attempt <= 10 {
		fmt.Printf("Request %d reached the API and FAILED!\n", attempt)
		return nil, errors.New("not implemented")
	} else {
		fmt.Printf("Request %d reached the API and SUCCEDED!\n", attempt)
		return []byte("success"), nil
	}

}
