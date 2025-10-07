package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hkumar1729/Url-shortener-API/internal/routes"
)

func main() {
	router := gin.Default()
	routes.RegisterRoutes(router)
	router.Run(":3000")
}
