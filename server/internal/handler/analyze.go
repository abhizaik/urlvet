package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/abhizaik/SafeSurf/internal/analyzer"
	"github.com/abhizaik/SafeSurf/internal/service/checks"
	"github.com/gin-gonic/gin"
)

// AnalyzeURLHandler runs a full safety analysis on the given URL.
//
//	@Summary		Full URL analysis
//	@Tags			Analysis
//	@Produce		json
//	@Param			url	query		string	true	"URL to analyse (max 2048 chars)"
//	@Success		200	{object}	analyzer.Response
//	@Failure		400	{object}	map[string]string
//	@Router			/analyze [get]
func AnalyzeURLHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url query param is required"})
		return
	}

	_, isValid, err := checks.IsValidURL(url)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"status": "ERROR", "error": "invalid url"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, _ := analyzer.Analyze(ctx, url)
	c.JSON(http.StatusOK, resp)
}
