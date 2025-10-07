package services

import (
	"context"
	"fmt"
	"time"

	"github.com/hkumar1729/Url-shortener-API/db"
	"github.com/hkumar1729/Url-shortener-API/internal/adapters/cache"
	"github.com/hkumar1729/Url-shortener-API/internal/adapters/database"
	"github.com/hkumar1729/Url-shortener-API/utils"
)

func CreateShortUrlService(c context.Context, OriginalUrl string, host string, scheme string) (string, error) {
	database.Init()
	defer database.PrismaClient.Prisma.Disconnect()

	newEntry, err := database.PrismaClient.URL.CreateOne(
		db.URL.OriginalURL.Set(OriginalUrl),
		db.URL.ShortenedURL.Set("temp"),
	).Exec(database.Ctx)

	if err != nil {
		return "", err
	}

	key := utils.GenerateUrlKey(newEntry.OriginalURL)

	_, err = database.PrismaClient.URL.FindUnique(
		db.URL.ID.Equals(newEntry.ID),
	).Update(
		db.URL.ShortenedURL.Set(key),
	).Exec(database.Ctx)

	if err != nil {
		return "", err
	}

	shortURL := fmt.Sprintf("%s://%s/%s", scheme, host, key)

	return shortURL, nil
}

func GetShortUrls() ([]db.URLModel, error) {
	database.Init()
	defer database.PrismaClient.Disconnect()

	urls, err := database.PrismaClient.URL.FindMany().Exec(database.Ctx)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func Redirect(key string) (*db.URLModel, error) {
	cache.InitRedis()
	defer cache.RedisClient.Close()

	cachedUrl, err := cache.RedisClient.Get(cache.Ctx, key).Result()
	if err == nil {
		return &db.URLModel{
			InnerURL: db.InnerURL{
				OriginalURL:  cachedUrl,
				ShortenedURL: &key,
			},
		}, nil
	}

	database.Init()
	defer database.PrismaClient.Prisma.Disconnect()

	url, err := database.PrismaClient.URL.FindFirst(
		db.URL.ShortenedURL.Equals(key),
	).Exec(database.Ctx)
	if err != nil {
		return nil, err
	}

	cache.RedisClient.Set(cache.Ctx, key, url.OriginalURL, 7*24*time.Hour)

	return url, nil
}
