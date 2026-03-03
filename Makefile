APP_NAME=repo-reviewer
MAIN_PATH=./cmd/api
BUILD_DIR=bin

# .PHONY: run build clean test tidy lint
.PHONY: all build run test clean watch docker-run docker-down itest tidy

all: build test

# Build the application
build:
	@echo "Building application..."

	@mkdir -p $(BUILD_DIR)

	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Run the application
run:
	@go run $(MAIN_PATH)

# Clean the binary
clean:
	@echo "Cleaning..."

	@rm -rf $(BUILD_DIR)

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Live reload
watch:
	@if command -v air > /dev/null; then\
						air; \
						echo "Watching application...";\
					else \
						read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice;\
						if ["$$choice" != "n"] && ["$$choice" != "N"]; then\
							go install github.com/air-verse/air@latest; \
							air;\
							echo "Watching application..."\
						else \
							echo "You chose not to install air. Existing...";\
							exit 1; \
						fi;\
					fi
# Tidy the application dependencies
tidy:
	go mod tidy

# Formatting application
lint:
	golangci-lint run

# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then\
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build;\
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down;
	fi