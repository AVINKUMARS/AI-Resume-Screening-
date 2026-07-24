package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scalingwolf/ai-resume-screening/backend/internal/handlers"
)

// New builds the Gin engine with CORS and all API routes registered.
func New(h *handlers.Handler, allowedOrigin string) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", h.Health)

	api := r.Group("/api")
	{
		api.POST("/resumes/upload", h.UploadResume)

		api.GET("/candidates", h.ListCandidates)
		api.GET("/candidates/:id", h.GetCandidate)
		api.DELETE("/candidates/:id", h.DeleteCandidate)

		api.POST("/jobs", h.CreateJob)
		api.GET("/jobs", h.ListJobs)
		api.GET("/jobs/:id", h.GetJob)
		api.DELETE("/jobs/:id", h.DeleteJob)
		api.POST("/jobs/:id/match", h.MatchJob)
	}

	return r
}
