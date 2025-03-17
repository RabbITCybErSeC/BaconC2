GO = go
BINARY_NAME = bacon-server
CLIENT_BINARY_NAME = bacon-client
SERV_SRC_DIR = ./server
CLIENT_SRC_DIR = ./client/
BIN_DIR = ./bin

LDFLAGS = -s -w
TAGS = go_sqlite,server
BUILD_FLAGS = -trimpath -tags "$(TAGS)" -ldflags "$(LDFLAGS)"

.PHONY: all clean compile-servers compile-client run-server run-client

all: compile-servers compile-client

$(BIN_DIR):
        @mkdir -p $(BIN_DIR)

compile-servers: $(BIN_DIR)
        @echo "Building server for darwin/arm64..."
        @cd $(SERV_SRC_DIR) && GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GO) build $(BUILD_FLAGS) -o ../$(BIN_DIR)/$(BINARY_NAME)_darwin-arm64
        @echo "Build completed."

compile-client: $(BIN_DIR)
        @echo "Building server for darwin/arm64..."
        @cd $(CLIENT_SRC_DIR) && GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GO) build $(BUILD_FLAGS) -o ../$(BIN_DIR)/$(CLIENT_BINARY_NAME)_darwin-arm64
        @echo "Build completed."

run-server:
        @echo "Running server..."
        @cd $(SERV_SRC_DIR) && $(GO) run -tags "$(TAGS)" .

run-client:
        @echo "Running client..."
        @cd $(CLIENT_SRC_DIR) && $(GO) run .

clean:
        @echo "Cleaning up..."
        @rm -rf $(BIN_DIR)
        @echo "Clean completed."


