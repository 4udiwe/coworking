FROM golang:1.21-alpine

WORKDIR /workspace/test

# Install curl for health checks and dependencies
RUN apk add --no-cache curl git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy test files
COPY . .

# Run tests
CMD ["go", "test", "-v", "-timeout", "120s", "./..."]
