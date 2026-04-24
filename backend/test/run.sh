#!/bin/bash

# Integration Test Runner Script
# This script starts all necessary services and runs the integration tests

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🧪 Coworking Integration Tests${NC}\n"

# Function to cleanup on exit
cleanup() {
    echo -e "${YELLOW}Cleaning up...${NC}"
    docker-compose -f docker-compose.test.yaml down -v || true
}

trap cleanup EXIT

# Start services
echo -e "${YELLOW}1. Starting services...${NC}"
docker-compose -f docker-compose.test.yaml up -d postgres_auth_test postgres_booking_test

echo -e "${YELLOW}   Waiting for databases...${NC}"
sleep 10

docker-compose -f docker-compose.test.yaml up -d auth-service booking-service

echo -e "${YELLOW}   Waiting for services to be ready...${NC}"
sleep 15

# Check health
echo -e "${YELLOW}2. Checking service health...${NC}"
max_retries=10
retry=0

while [ $retry -lt $max_retries ]; do
    auth_health=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8081/health || echo "000")
    booking_health=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8082/health || echo "000")
    
    if [ "$auth_health" = "200" ] && [ "$booking_health" = "200" ]; then
        echo -e "${GREEN}   ✓ Auth Service: OK${NC}"
        echo -e "${GREEN}   ✓ Booking Service: OK${NC}"
        break
    fi
    
    echo -e "${YELLOW}   Retrying... (attempt $((retry+1))/$max_retries)${NC}"
    sleep 3
    retry=$((retry+1))
done

if [ "$auth_health" != "200" ] || [ "$booking_health" != "200" ]; then
    echo -e "${RED}✗ Services failed to start${NC}"
    echo -e "${RED}   Auth Service: $auth_health${NC}"
    echo -e "${RED}   Booking Service: $booking_health${NC}"
    exit 1
fi

# Run tests
echo -e "${YELLOW}3. Running tests...${NC}"
docker-compose -f docker-compose.test.yaml run --rm test

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Tests failed!${NC}"
    exit 1
fi
