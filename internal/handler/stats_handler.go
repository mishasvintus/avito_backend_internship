package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mishasvintus/avito_backend_internship/internal/service"
)

// StatsHandler handles statistics HTTP requests.
type StatsHandler struct {
	statsService StatsServiceInterface
}

// StatsServiceInterface defines the interface for statistics operations.
type StatsServiceInterface interface {
	GetStatistics() (*service.Statistics, error)
}

// NewStatsHandler creates a new stats handler.
func NewStatsHandler(statsService StatsServiceInterface) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// GetStatistics handles GET /stats.
func (h *StatsHandler) GetStatistics(c *gin.Context) {
	stats, err := h.statsService.GetStatistics()
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	response := StatisticsResponse{
		Overall: struct {
			TotalPRs         int64 `json:"total_prs"`
			TotalAssignments int64 `json:"total_assignments"`
			TotalUsers       int64 `json:"total_users"`
			TotalTeams       int64 `json:"total_teams"`
		}{
			TotalPRs:         stats.Overall.TotalPRs,
			TotalAssignments: stats.Overall.TotalAssignments,
			TotalUsers:       stats.Overall.TotalUsers,
			TotalTeams:       stats.Overall.TotalTeams,
		},
		ReviewerStats: make([]ReviewerStatResponse, len(stats.ReviewerStats)),
		AuthorStats:   make([]AuthorStatResponse, len(stats.AuthorStats)),
	}

	for i, rs := range stats.ReviewerStats {
		response.ReviewerStats[i] = ReviewerStatResponse{
			UserID:   rs.UserID,
			Username: rs.Username,
			Count:    rs.Count,
		}
	}

	for i, as := range stats.AuthorStats {
		response.AuthorStats[i] = AuthorStatResponse{
			UserID:   as.UserID,
			Username: as.Username,
			Count:    as.Count,
		}
	}

	c.JSON(http.StatusOK, response)
}
