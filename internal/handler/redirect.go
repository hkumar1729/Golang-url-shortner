package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hkumar1729/Url-shortener-API/internal/core/services"
)

func RedirectUrl(c *gin.Context) {
	key := c.Param("shorturl")

	url, err := services.Redirect(key)
	if err != nil || url == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "short URL not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}
