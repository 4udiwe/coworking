# Step 1: Modules caching
FROM golang:1.25-alpine AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.25-alpine AS builder

# Устанавливаем build dependencies для vips
RUN apk add --no-cache \
    build-base \
    vips-dev \
    vips-tools \
    pkgconfig

COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=1 GOOS=linux go build -tags vips -o /app/media-service ./cmd/main.go

# Step 3: Final
FROM alpine:3.22

# Устанавливаем runtime зависимости для vips
RUN apk add --no-cache \
    vips \
    vips-tools \
    ca-certificates \
    curl

COPY --from=builder /app/config /config
COPY --from=builder /app/media-service /app/media-service

WORKDIR /app
EXPOSE 8084
CMD ["/app/media-service"]