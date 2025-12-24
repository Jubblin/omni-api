.PHONY: test build-gui run-gui clean tidy version version-patch version-minor version-major

# Run all tests
test:
	go test ./...

# Build the GUI application
build-gui:
	@VERSION=$$(cat VERSION 2>/dev/null || echo "dev"); \
	CGO_LDFLAGS="-Wl,-no_warn_duplicate_libraries" go build -ldflags "-X github.com/jubblin/omni-api/main.Version=$$VERSION" -o omni-gui .

# Run the GUI application
run-gui:
	go run .

# Tidy go modules
tidy:
	go mod tidy

# Clean up
clean:
	rm -f omni-api omni-gui
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

