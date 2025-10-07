package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hkumar1729/Url-shortener-API/db"
	"github.com/hkumar1729/Url-shortener-API/entity"
	"github.com/hkumar1729/Url-shortener-API/internal/adapters/cache"
	"github.com/hkumar1729/Url-shortener-API/internal/database"
	"github.com/hkumar1729/Url-shortener-API/utils"
)

func getShortUrls(c *gin.Context) {
	database.Init()
	defer database.PrismaClient.Prisma.Disconnect()
	urls, err := database.PrismaClient.URL.FindMany().Exec(database.Ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, urls)
}

func createShortUrl(c *gin.Context) {
	var input entity.Url
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	database.Init()
	defer database.PrismaClient.Prisma.Disconnect()

	newEntry, err := database.PrismaClient.URL.CreateOne(
		db.URL.OriginalURL.Set(input.Url),
		db.URL.ShortenedURL.Set("temp"),
	).Exec(database.Ctx)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errro": err.Error()})
	}

	key := utils.GenerateUrlKey(newEntry.OriginalURL)

	_, err = database.PrismaClient.URL.FindUnique(
		db.URL.ID.Equals(newEntry.ID),
	).Update(
		db.URL.ShortenedURL.Set(key),
	).Exec(database.Ctx)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	host := c.Request.Host // domain name or localhost

	shortURL := fmt.Sprintf("%s://%s/%s", scheme, host, key)

	c.JSON(http.StatusOK, gin.H{
		"short-url": shortURL,
	})

}

func redirectUrl(c *gin.Context) {
	key := c.Param("shorturl")

	cache.InitRedis()
	defer cache.RedisClient.Close()

	cachedUrl, err := cache.RedisClient.Get(cache.Ctx, key).Result()
	if err == nil {
		c.Redirect(http.StatusMovedPermanently, cachedUrl)
		return
	}

	database.Init()
	defer database.PrismaClient.Prisma.Disconnect()

	url, err := database.PrismaClient.URL.FindFirst(
		db.URL.ShortenedURL.Equals(key),
	).Exec(database.Ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cache.RedisClient.Set(cache.Ctx, key, url.OriginalURL, 7*24*time.Hour)

	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

func main() {

	server := gin.Default()

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	server.GET("/shorturl", getShortUrls)
	server.POST("/shorturl", createShortUrl)
	server.GET("/:shorturl", redirectUrl)

	server.Run(":3000")
}
