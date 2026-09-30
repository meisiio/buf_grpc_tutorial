package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// This represents ANY function we want to run (like making an HTTP call)
type ActionFunc func() error

func doWithRetry(maxRetries int, baseDelay time.Duration, action ActionFunc) error {
	delay := baseDelay

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// 1. Try executing the action!
		err := action()
		if err == nil {
			return nil // Success! We don't need to retry.
		}

		fmt.Printf("Attempt %d failed. Retrying...\n", attempt)

		// 2. TODO: Calculate Jitter!
		// Generate a random duration between 0 and 500 Milliseconds using rand.Intn(500)
		jitter := time.Duration(rand.Intn(500)) * time.Millisecond
		// 3. TODO: Add the jitter to the current `delay`

		fmt.Printf("Waiting for %v before next attempt...\n", delay+jitter)

		// 4. TODO: Sleep for the total time (delay + jitter)
		time.Sleep(delay + jitter)
		// 5. TODO: Exponential Backoff! Multiply the base `delay` by 2 for the NEXT iteration.
		delay *= 2
	}

	return errors.New("all retries failed")
}

func main() {
	// A fake API that always fails
	flakyAction := func() error {
		return errors.New("503 Service Unavailable")
	}

	fmt.Println("Starting API Call...")

	err := doWithRetry(4, 1*time.Second, flakyAction)
	if err != nil {
		fmt.Println("Final result:", err)
	}
}
