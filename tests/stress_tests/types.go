package stress_tests

import (
	"sync"
	"time"
)

const (
	teamName = "loadtest"
	baseURL  = "http://localhost:8080"
	duration = 30 * time.Second
)

// Result represents a single request result.
type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
	Endpoint   string
}

var (
	// reassignCounter is an atomic counter for round-robin PR selection.
	reassignCounter int64

	// currentReviewers stores current reviewers for each PR.
	currentReviewers sync.Map

	// reviewerMutex protects reviewer map updates.
	reviewerMutex sync.Mutex
)
