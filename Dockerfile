FROM golang:1.25 AS builder

WORKDIR /api

COPY go.mod go.sum ./
COPY entrypoint.sh .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping ./cmd/server/main.go

EXPOSE 3000

ENTRYPOINT ["./entrypoint.sh"]