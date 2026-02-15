cover:
	go test -count=1 -coverprofile=coverage.out $(shell go list ./internal/... | grep -E "internal/service" | grep -v -E "mocks")
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out