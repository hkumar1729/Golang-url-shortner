package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hkumar1729/Url-shortener-API/internal/handler"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/", handler.HealthCheck)
	router.GET("/shorturl", handler.GetShortUrls)
	router.POST("/shorturl", handler.CreateShortUrl)
	router.GET("/:shorturl", handler.RedirectUrl)
}
