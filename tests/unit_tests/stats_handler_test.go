package unit_tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mishasvintus/avito_backend_internship/internal/handler"
	"github.com/mishasvintus/avito_backend_internship/internal/repository/stats"
	"github.com/mishasvintus/avito_backend_internship/internal/service"
	handlermocks "github.com/mishasvintus/avito_backend_internship/tests/mocks"
)

func TestStatsHandler_GetStatistics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		mockSetup        func(*handlermocks.MockStatsServiceInterface)
		expectedStatus   int
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success - returns statistics",
			mockSetup: func(m *handlermocks.MockStatsServiceInterface) {
				m.EXPECT().GetStatistics().Return(&service.Statistics{
					Overall: &stats.OverallStats{
						TotalPRs:         5,
						TotalAssignments: 10,
						TotalUsers:       3,
						TotalTeams:       2,
					},
					ReviewerStats: []stats.ReviewerStat{
						{
							UserID:   "user1",
							Username: "reviewer1",
							Count:    5,
						},
						{
							UserID:   "user2",
							Username: "reviewer2",
							Count:    3,
						},
					},
					AuthorStats: []stats.AuthorStat{
						{
							UserID:   "user1",
							Username: "author1",
							Count:    3,
						},
						{
							UserID:   "user2",
							Username: "author2",
							Count:    2,
						},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response handler.StatisticsResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Check overall stats
				assert.Equal(t, int64(5), response.Overall.TotalPRs)
				assert.Equal(t, int64(10), response.Overall.TotalAssignments)
				assert.Equal(t, int64(3), response.Overall.TotalUsers)
				assert.Equal(t, int64(2), response.Overall.TotalTeams)

				// Check reviewer stats
				assert.Len(t, response.ReviewerStats, 2)
				assert.Equal(t, "user1", response.ReviewerStats[0].UserID)
				assert.Equal(t, "reviewer1", response.ReviewerStats[0].Username)
				assert.Equal(t, int64(5), response.ReviewerStats[0].Count)
				assert.Equal(t, "user2", response.ReviewerStats[1].UserID)
				assert.Equal(t, "reviewer2", response.ReviewerStats[1].Username)
				assert.Equal(t, int64(3), response.ReviewerStats[1].Count)

				// Check author stats
				assert.Len(t, response.AuthorStats, 2)
				assert.Equal(t, "user1", response.AuthorStats[0].UserID)
				assert.Equal(t, "author1", response.AuthorStats[0].Username)
				assert.Equal(t, int64(3), response.AuthorStats[0].Count)
				assert.Equal(t, "user2", response.AuthorStats[1].UserID)
				assert.Equal(t, "author2", response.AuthorStats[1].Username)
				assert.Equal(t, int64(2), response.AuthorStats[1].Count)
			},
		},
		{
			name: "success - empty statistics",
			mockSetup: func(m *handlermocks.MockStatsServiceInterface) {
				m.EXPECT().GetStatistics().Return(&service.Statistics{
					Overall: &stats.OverallStats{
						TotalPRs:         0,
						TotalAssignments: 0,
						TotalUsers:       0,
						TotalTeams:       0,
					},
					ReviewerStats: []stats.ReviewerStat{},
					AuthorStats:   []stats.AuthorStat{},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response handler.StatisticsResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, int64(0), response.Overall.TotalPRs)
				assert.Equal(t, int64(0), response.Overall.TotalAssignments)
				assert.Equal(t, int64(0), response.Overall.TotalUsers)
				assert.Equal(t, int64(0), response.Overall.TotalTeams)
				assert.Empty(t, response.ReviewerStats)
				assert.Empty(t, response.AuthorStats)
			},
		},
		{
			name: "error - internal error from service",
			mockSetup: func(m *handlermocks.MockStatsServiceInterface) {
				m.EXPECT().GetStatistics().Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response handler.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response.Error.Message, "assert.AnError")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := handlermocks.NewMockStatsServiceInterface(t)
			tt.mockSetup(mockService)

			statsHandler := handler.NewStatsHandler(mockService)

			req, err := http.NewRequest(http.MethodGet, "/stats", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			statsHandler.GetStatistics(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.validateResponse(t, w)
		})
	}
}
