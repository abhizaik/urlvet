package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler returns the service liveness status.
//
//	@Summary		Health check
//	@Tags			Utility
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
