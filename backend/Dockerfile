FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/bin/api ./cmd/api

RUN wget -O /migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz && \
    tar -xzf /migrate.tar.gz && \
    mv migrate /app/bin/migrate

FROM alpine:latest

WORKDIR /app
RUN apk add --no-cache poppler-utils

COPY --from=builder /app/bin/api .
COPY --from=builder /app/bin/migrate .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

EXPOSE 8090
CMD ["./api"]