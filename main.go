package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hkumar1729/Url-shortener-API/db"
	"github.com/hkumar1729/Url-shortener-API/internal/database"
)

func main() {

	server := gin.Default()

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	database.Init()
	defer database.PrismaClient.Prisma.Disconnect()

	newURL, err := database.PrismaClient.URL.CreateOne(
		db.URL.OriginalURL.Set("https://example.com"),
		db.URL.ShortenedURL.Set("exmpl"),
	).Exec(database.Ctx)

	if err != nil {
		panic(err)
	}

	fmt.Println("Created URL ID:", newURL.ID)

	server.Run(":3000")
}
