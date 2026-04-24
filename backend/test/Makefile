.PHONY: help test test-up test-down test-logs test-clean test-rebuild push-test

help:
	@echo "Integration Test Commands:"
	@echo "  make test          - Run integration tests (starts services, runs tests, cleans up)"
	@echo "  make test-up       - Start test services (auth, booking, postgres)"
	@echo "  make test-down     - Stop and remove all test containers"
	@echo "  make test-logs     - View logs from running test services"
	@echo "  make test-clean    - Clean up test volumes and networks"
	@echo "  make test-rebuild  - Rebuild test containers and run tests"
	@echo "  make build-test    - Build test services"

# Run full test flow: start services → run tests → clean up
test: test-down test-up test-run test-down
	@echo "✅ Integration tests completed!"

# Start only the services without running tests
test-up:
	@echo "Starting test services..."
	docker-compose -f docker-compose.test.yaml up -d postgres_auth_test postgres_booking_test
	@echo "Waiting for databases to be ready..."
	sleep 5
	docker-compose -f docker-compose.test.yaml up -d auth-service booking-service
	@echo "Waiting for services to be ready..."
	sleep 5
	@echo "✅ Services started. Check status with: make test-logs"

# Stop services
test-down:
	@echo "Stopping test services..."
	docker-compose -f docker-compose.test.yaml down -v || true
	@echo "✅ Services stopped and volumes removed"

# View logs
test-logs:
	docker-compose -f docker-compose.test.yaml logs -f

# Run the actual tests
test-run:
	@echo "Running integration tests..."
	docker-compose -f docker-compose.test.yaml run --rm test

# Clean everything
test-clean:
	docker-compose -f docker-compose.test.yaml down -v --remove-orphans
	docker network prune -f || true
	@echo "✅ Cleanup completed"

# Rebuild and run
test-rebuild: test-down
	docker-compose -f docker-compose.test.yaml build --no-cache
	$(MAKE) test

# Just build the test images
build-test:
	docker-compose -f docker-compose.test.yaml build

# Local testing (requires services running locally)
test-local:
	@echo "Running tests against localhost services..."
	cd $(CURDIR) && go test -v -timeout 120s ./...

# Check test service health
test-health:
	@echo "Auth Service: $$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8081/health || echo 'DOWN')"
	@echo "Booking Service: $$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8082/health || echo 'DOWN')"
	@echo "PostgreSQL Auth: $$(docker exec postgres_auth_test pg_isready -U authuser -d authdb 2>/dev/null || echo 'DOWN')"
	@echo "PostgreSQL Booking: $$(docker exec postgres_booking_test pg_isready -U bookinguser -d bookingdb 2>/dev/null || echo 'DOWN')"
