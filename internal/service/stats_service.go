package service

import (
	"database/sql"
	"fmt"

	"github.com/mishasvintus/avito_backend_internship/internal/repository/stats"
)

// StatsService handles statistics business logic.
type StatsService struct {
	db *sql.DB
}

// NewStatsService creates a new stats service.
func NewStatsService(db *sql.DB) *StatsService {
	return &StatsService{db: db}
}

// Statistics represents all statistics.
type Statistics struct {
	Overall       *stats.OverallStats
	ReviewerStats []stats.ReviewerStat
	AuthorStats   []stats.AuthorStat
}

// GetStatistics returns all statistics.
func (s *StatsService) GetStatistics() (*Statistics, error) {
	overall, err := stats.GetOverallStats(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get overall stats: %w", err)
	}

	reviewerStats, err := stats.GetReviewerStats(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewer stats: %w", err)
	}

	authorStats, err := stats.GetAuthorStats(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get author stats: %w", err)
	}

	return &Statistics{
		Overall:       overall,
		ReviewerStats: reviewerStats,
		AuthorStats:   authorStats,
	}, nil
}
