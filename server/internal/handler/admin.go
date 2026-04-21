package handler

import (
	"net/http"

	"github.com/abhizaik/SafeSurf/internal/store"
	"github.com/gin-gonic/gin"
)

// AdminStatsHandler returns live scan metrics derived from the in-memory ring buffer.
func AdminStatsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, store.GetStats())
}

// AdminRecentHandler returns the last N scan records (newest first).
func AdminRecentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"scans": store.RecentScans()})
}

// AdminErrorsHandler returns the last N error/panic records (newest first).
func AdminErrorsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"errors": store.RecentErrors()})
}
