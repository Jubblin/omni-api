# Integration Test Suite

> 📖 [Back to README](../README.md)

This directory contains integration tests that test the entire API against a configured Omni host.

## Overview

The integration test suite validates that all API endpoints work correctly against a real Omni instance. Unlike unit tests which use mocks, these tests require:

- A running Omni instance
- Valid authentication credentials
- A running API server (or the tests will start one)

## Prerequisites

1. **Omni Instance**: You must have access to a running Omni instance
2. **Authentication**: Configure authentication via environment variables
3. **API Server**: The API server must be running (start with `make run` in a separate terminal)

## Configuration

Integration tests use the same environment variables as the main application:

### Required

- `OMNI_ENDPOINT` - The Omni API endpoint URL
- `INTEGRATION_TESTS=true` - Must be set to enable integration tests

### Authentication (choose one)

**Service Account:**

- `OMNI_SERVICE_ACCOUNT` or `OMNI_SERVICE_ACCOUNT_KEY`

**PGP Authentication:**

- `OMNI_CONTEXT`
- `OMNI_IDENTITY`
- `OMNI_KEYS_DIR` (optional)

### Optional

- `TEST_PORT` - Port for API server being tested (default: 8080, or set `PORT` env var when starting server)
- `OMNI_INSECURE=true` - Skip TLS verification

## Running Integration Tests

### Prerequisites Check

Integration tests are skipped by default. To run them:

```bash
export INTEGRATION_TESTS=true
export OMNI_ENDPOINT="https://your-omni-instance.com"
export OMNI_SERVICE_ACCOUNT_KEY="your-service-account-key"
```

### Run All Integration Tests

```bash
go test ./integration/... -v
```

### Run Specific Test

```bash
go test ./integration/... -v -run TestClustersEndpoints
```

### Run with Server

Start the API server in one terminal, then run tests in another:

```bash
# Terminal 1: Start the API server
export OMNI_ENDPOINT="https://your-omni-instance.com"
export OMNI_SERVICE_ACCOUNT_KEY="your-key"
make run

# Terminal 2: Run integration tests
export INTEGRATION_TESTS=true
export OMNI_ENDPOINT="https://your-omni-instance.com"
export OMNI_SERVICE_ACCOUNT_KEY="your-key"
go test ./integration/... -v
```

**Note**: The API server must be running before tests execute. Tests will fail if the server is not accessible.

## Test Coverage

The integration test suite covers:

- ✅ Health checks (root redirect, Swagger UI)
- ✅ Cluster endpoints (list, get, status, metrics, bootstrap, etc.)
- ✅ Machine endpoints (list, get, labels, extensions, metrics, etc.)
- ✅ Machine set endpoints (list, get, status, destroy status)
- ✅ Cluster machine endpoints (list, get, status, config, etc.)
- ✅ Other resource endpoints (config patches, machine classes, etc.)
- ✅ Query parameter filtering
- ✅ Link generation (_links in responses)
- ✅ Error handling (404 responses)

## Test Structure

- `api_test.go` - Main integration test suite
- Tests are organized by resource type
- Tests handle missing resources gracefully (404 is acceptable)
- Tests log useful information for debugging

## Best Practices

1. **Use Test Environment**: Run integration tests against a test/staging Omni instance, not production
2. **Clean State**: Tests assume resources may or may not exist
3. **Idempotent**: Tests don't modify resources, only read them
4. **Graceful Degradation**: Tests handle missing resources (404) as acceptable

## Troubleshooting

### Tests Skipped

If tests are skipped, ensure `INTEGRATION_TESTS=true` is set.

### Connection Errors

- Verify `OMNI_ENDPOINT` is correct and accessible
- Check authentication credentials are valid
- Ensure network connectivity to Omni instance

### 404 Errors

404 errors are expected for resources that don't exist. Tests handle this gracefully.

### Timeout Errors

If tests timeout:

- Check Omni instance is responsive
- Increase timeout in test client if needed
- Verify network connectivity

## Related Documentation

- 📖 [README.md](../README.md) - Main project documentation
- 📋 [RESOURCES.md](../RESOURCES.md) - Available Omni resources
- 📊 [TEST_COVERAGE.md](../TEST_COVERAGE.md) - Unit test coverage report
- 🔧 [MACHINE_ENHANCEMENTS.md](../MACHINE_ENHANCEMENTS.md) - Machine endpoint enhancements
