# Step 1: Modules caching
FROM golang:1.25 AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.25 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o /app/gateway ./cmd/main.go

# Step 3: Final
FROM alpine:3.22

COPY --from=builder /app/config /config
COPY --from=builder /app/gateway /app/gateway

WORKDIR /app
EXPOSE 8080
CMD ["/app/gateway"]