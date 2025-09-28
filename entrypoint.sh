#!/bin/sh
set -e

echo "Applying Prisma schema..."
go run github.com/steebchen/prisma-client-go db push

echo "Starting Go server..."
/docker-gs-ping