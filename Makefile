APP_NAME=repo-reviewer
MAIN_PATH=./cmd/server
BUILD_DIR=bin

.PHONY: run build clean test tidy lint

run:
	go run $(MAIN_PATH)

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...

tidy:
	go mod tidy

lint:
	golangci-lint run