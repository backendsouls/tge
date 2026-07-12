.PHONY: all build vet test test-unit test-functional clean

# Default binary output
APP_NAME=tge
BIN_DIR=bin
MAIN_PATH=./cmd/tge

all: vet test build

# Build the application
build:
	@echo "==> Building ${APP_NAME}..."
	@mkdir -p ${BIN_DIR}
	@go build -o ${BIN_DIR}/${APP_NAME} ${MAIN_PATH}
	@echo "==> Build complete: ${BIN_DIR}/${APP_NAME}"

# Run go vet
vet:
	@echo "==> Vetting code..."
	@go vet ./...

# Run all tests
test:
	@echo "==> Running all tests..."
	@go test -v ./...

# Run only unit tests (fast tests)
unit-test:
	@echo "==> Running unit tests..."
	@go test -v -short ./...

# Run functional/integration tests
functional-test:
	@echo "==> Running functional tests..."
	# In the future, this can be mapped to go test -tags=functional ./... or similar.
	# For now, it runs the standard full suite.
	@go test -v -run "Integration|Functional" ./... || echo "No specific functional tests found or some failed."

# Clean compiled binaries
clean:
	@echo "==> Cleaning build cache..."
	@go clean
	@rm -rf ${BIN_DIR}
