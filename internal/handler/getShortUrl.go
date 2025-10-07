package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hkumar1729/Url-shortener-API/internal/core/services"
)

func GetShortUrls(c *gin.Context) {
	urls, err := services.GetShortUrls()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, urls)
}
