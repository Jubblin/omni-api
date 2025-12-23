.PHONY: test build run swagger clean tidy version version-patch version-minor version-major

# Run all tests
test:
	go test ./...

# Run integration tests (requires INTEGRATION_TESTS=true and configured Omni host)
test-integration:
	@if [ "$(INTEGRATION_TESTS)" != "true" ]; then \
		echo "Integration tests require INTEGRATION_TESTS=true"; \
		echo "Example: INTEGRATION_TESTS=true OMNI_ENDPOINT=https://omni.example.com OMNI_SERVICE_ACCOUNT_KEY=... make test-integration"; \
		exit 1; \
	fi
	go test ./integration/... -v

# Build the binary
build:
	@VERSION=$$(cat VERSION 2>/dev/null || echo "dev"); \
	go build -ldflags "-X github.com/jubblin/omni-api/internal/api/handlers.Version=$$VERSION -X github.com/jubblin/omni-api/main.Version=$$VERSION" -o omni-api main.go

# Run the application
run:
	go run main.go

# Generate Swagger documentation
swagger:
	go run github.com/swaggo/swag/cmd/swag init

# Tidy go modules
tidy:
	go mod tidy

# Clean up
clean:
	rm -f omni-api
	rm -rf docs/

# Version management
version:
	@./scripts/version.sh get

version-patch:
	@./scripts/version.sh patch

version-minor:
	@./scripts/version.sh minor

version-major:
	@./scripts/version.sh major

