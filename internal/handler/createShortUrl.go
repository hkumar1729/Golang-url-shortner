package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hkumar1729/Url-shortener-API/entity"
	"github.com/hkumar1729/Url-shortener-API/internal/core/services"
)

func CreateShortUrl(c *gin.Context) {
	var input entity.Url
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	host := c.Request.Host // domain name or localhost

	shortUrl, err := services.CreateShortUrl(c, input.Url, host, scheme)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"short-url": shortUrl,
	})

}
