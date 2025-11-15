package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mishasvintus/avito_backend_internship/internal/handler"
)

// SetupRoutes configures all API routes.
func SetupRoutes(
	teamHandler *handler.TeamHandler,
	userHandler *handler.UserHandler,
	prHandler *handler.PRHandler,
) *gin.Engine {
	r := gin.Default()

	// Team endpoints
	r.POST("/team/add", teamHandler.AddTeam)
	r.GET("/team/get", teamHandler.GetTeam)

	// User endpoints
	r.POST("/users/setIsActive", userHandler.SetIsActive)
	r.GET("/users/getReview", userHandler.GetReview)

	// Pull Request endpoints
	r.POST("/pullRequest/create", prHandler.CreatePR)
	r.POST("/pullRequest/merge", prHandler.MergePR)
	r.POST("/pullRequest/reassign", prHandler.ReassignPR)

	return r
}
