package database

import (
	"context"

	prisma "github.com/hkumar1729/Url-shortener-API/db"
)

var PrismaClient *prisma.PrismaClient
var Ctx = context.Background()

func Init() {
	PrismaClient = prisma.NewClient()
	if err := PrismaClient.Prisma.Connect(); err != nil {
		panic(err)
	}
}
