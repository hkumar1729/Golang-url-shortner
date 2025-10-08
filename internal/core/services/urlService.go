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

func CreateShortUrl(c context.Context, OriginalUrl string, host string, scheme string) (string, error) {
	database.Init()
	defer database.PrismaClient.Prisma.Disconnect()

	existing, err := database.PrismaClient.URL.FindFirst(
		db.URL.OriginalURL.Equals(OriginalUrl),
	).Exec(database.Ctx)
	if err != nil && err != db.ErrNotFound {
		return "", err
	}

	if existing != nil {
		shortURL := fmt.Sprintf("%s://%s/%s", scheme, host, existing.ShortenedURL)
		return shortURL, nil
	}

	key := utils.GenerateUrlKey(OriginalUrl)

	_, err = database.PrismaClient.URL.CreateOne(
		db.URL.OriginalURL.Set(OriginalUrl),
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

	database.Init()
	defer database.PrismaClient.Prisma.Disconnect()

	cachedUrl, err := cache.RedisClient.Get(cache.Ctx, key).Result()
	if err == nil {

		_, err := database.PrismaClient.URL.FindUnique(
			db.URL.ShortenedURL.Equals(key),
		).Update(
			db.URL.CountClick.Increment(1),
		).Exec(database.Ctx)
		if err != nil {
			return nil, err
		}

		return &db.URLModel{
			InnerURL: db.InnerURL{
				OriginalURL:  cachedUrl,
				ShortenedURL: key,
			},
		}, nil
	}

	url, err := database.PrismaClient.URL.FindUnique(
		db.URL.ShortenedURL.Equals(key),
	).Update(
		db.URL.CountClick.Increment(1),
	).Exec(database.Ctx)
	if err != nil {
		return nil, err
	}

	cache.RedisClient.Set(cache.Ctx, key, url.OriginalURL, 7*24*time.Hour)

	return url, nil
}
