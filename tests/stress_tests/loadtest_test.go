package stress_tests

import (
	"sync"
	"testing"
	"time"
)

const (
	targetRPS = 5
)

// TestLoadTest performs uniform load testing on multiple endpoints.
func TestLoadTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	t.Log("Setting up test data...")
	setupTestData(t)

	t.Logf("Starting load test: %d RPS for %v\n", targetRPS, duration)

	endpoints := []struct {
		name string
		fn   func(chan<- Result)
	}{
		{"GET /team/get", testGetTeam},
		{"GET /users/getReview", testGetUserReviews},
		{"POST /users/setIsActive", testSetIsActive},
	}

	var allResults []Result
	var mu sync.Mutex

	// Test all endpoints in parallel
	results := make(chan Result, 10000)
	var wg sync.WaitGroup

	done := make(chan bool)
	for _, ep := range endpoints {
		t.Logf("Testing %s...", ep.name)
		interval := time.Second / targetRPS
		ticker := time.NewTicker(interval)

		go func(endpoint struct {
			name string
			fn   func(chan<- Result)
		}) {
			for {
				select {
				case <-ticker.C:
					wg.Add(1)
					go func() {
						defer wg.Done()
						endpoint.fn(results)
					}()
				case <-done:
					ticker.Stop()
					return
				}
			}
		}(ep)
	}

	// Collector
	var collectorWg sync.WaitGroup
	collectorWg.Add(1)
	go func() {
		defer collectorWg.Done()
		for r := range results {
			mu.Lock()
			allResults = append(allResults, r)
			mu.Unlock()
		}
	}()

	time.Sleep(duration)
	close(done)
	wg.Wait()
	close(results)
	collectorWg.Wait()

	t.Log("\n=== Load Test Results ===")
	analyzeResults(t, allResults, duration)
}
