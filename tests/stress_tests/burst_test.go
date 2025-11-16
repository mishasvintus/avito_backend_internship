package stress_tests

import (
	"sync"
	"testing"
	"time"
)

const (
	burstSize     = 5
	burstInterval = 1 * time.Second
)

// TestBurstLoadTest performs burst load testing on multiple endpoints.
func TestBurstLoadTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping burst load test in short mode")
	}

	t.Log("Setting up test data...")
	setupTestData(t)

	t.Logf("Starting burst load test: %d requests per burst, burst every %v, for %v\n", burstSize, burstInterval, duration)

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
		t.Logf("Testing %s with bursts...", ep.name)
		burstTicker := time.NewTicker(burstInterval)

		go func(endpoint struct {
			name string
			fn   func(chan<- Result)
		}) {
			for {
				select {
				case <-burstTicker.C:
					for i := 0; i < burstSize; i++ {
						wg.Add(1)
						go func() {
							defer wg.Done()
							endpoint.fn(results)
						}()
					}
				case <-done:
					burstTicker.Stop()
					return
				}
			}
		}(ep)
	}

	// Collect results
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

	t.Log("\n=== Burst Load Test Results ===")
	analyzeResults(t, allResults, duration)
}
