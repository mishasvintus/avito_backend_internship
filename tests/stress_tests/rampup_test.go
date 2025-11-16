package stress_tests

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestRampUpReassignPR performs ramp-up load testing on PR reassign endpoint.
func TestRampUpReassignPR(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ramp-up reassign test in short mode")
	}

	atomic.StoreInt64(&reassignCounter, 0)
	currentReviewers = sync.Map{}

	t.Log("Setting up test data for reassign test...")
	prIDs, reviewerIDs := setupReassignTestData(t)

	t.Log("Starting ramp-up load test for PR reassign:")
	t.Log("  Initial RPS: 5")
	t.Log("  Ramp-up: +10 RPS every 10 seconds")
	t.Logf("  Duration: %v", duration)

	results := make(chan Result, 10000)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allResults []Result

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

	rampUpInterval := 10 * time.Second
	done := make(chan bool)

	var currentRPS int64 = 5

	go func() {
		startTime := time.Now()
		rampUpTicker := time.NewTicker(rampUpInterval)
		defer rampUpTicker.Stop()

		requestTicker := time.NewTicker(10 * time.Millisecond)
		defer requestTicker.Stop()

		lastRequestTime := startTime

		go func() {
			for {
				select {
				case <-done:
					return
				case <-requestTicker.C:
					now := time.Now()
					rps := atomic.LoadInt64(&currentRPS)
					interval := time.Second / time.Duration(rps)

					// Send request if enough time has passed
					if now.Sub(lastRequestTime) >= interval {
						wg.Add(1)
						go func() {
							defer wg.Done()
							testReassignPR(results, prIDs, reviewerIDs)
						}()
						lastRequestTime = now
					}
				}
			}
		}()

		// Ramp up every 10 seconds
		for {
			select {
			case <-done:
				return
			case <-rampUpTicker.C:
				elapsed := time.Since(startTime)
				if elapsed >= duration {
					return
				}
				newRPS := atomic.AddInt64(&currentRPS, 10)
				t.Logf("Ramping up to %d RPS at %v", newRPS, elapsed.Round(time.Second))
			}
		}
	}()

	time.Sleep(duration)
	close(done)
	wg.Wait()
	close(results)
	collectorWg.Wait()

	t.Log("\n=== Ramp-Up Reassign PR Test Results ===")
	analyzeResultsWithPercentiles(t, allResults, duration)
}
