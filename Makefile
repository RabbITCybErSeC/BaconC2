GO = go
BINARY_NAME = bacon-server
CLIENT_BINARY_NAME = bacon-client
SERV_SRC_DIR = ./server
CLIENT_SRC_DIR = ./client
FRONT_END_DIR = ./becongui
BIN_DIR = ./bin

TARGET_OS ?= darwin
TARGET_ARCH ?= arm64
TARGET = $(TARGET_OS)-$(TARGET_ARCH)

LDFLAGS = -s -w
SERVER_TAGS = go_sqlite,server
CLIENT_TAGS =
SERVER_BUILD_FLAGS = -trimpath -tags "$(SERVER_TAGS)" -ldflags "$(LDFLAGS)"
CLIENT_BUILD_FLAGS = -trimpath -tags "$(CLIENT_TAGS)" -ldflags "$(LDFLAGS)"

SERVER_OUTPUT = $(BIN_DIR)/$(BINARY_NAME)_$(TARGET)
CLIENT_OUTPUT = $(BIN_DIR)/$(CLIENT_BINARY_NAME)_$(TARGET)

.PHONY: all clean compile-server compile-client run-server run-client build-front-end

all: compile-server compile-client

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

build-front-end: 
	@echo "Building front-end"
	@cd $(FRONT_END_DIR) && npm run build
	
compile-server: $(BIN_DIR)
	@echo "Building server for $(TARGET)..."
	@cd $(SERV_SRC_DIR) && GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) CGO_ENABLED=1 $(GO) build $(SERVER_BUILD_FLAGS) -o ../$(SERVER_OUTPUT) .
	@echo "Server build completed: $(SERVER_OUTPUT)"

compile-client: $(BIN_DIR)
	@echo "Building client for $(TARGET)..."
	@cd $(CLIENT_SRC_DIR) && GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) CGO_ENABLED=0 $(GO) build $(CLIENT_BUILD_FLAGS) -o ../$(CLIENT_OUTPUT) .
	@echo "Client build completed: $(CLIENT_OUTPUT)"

run-server: build-front-end
	@echo "Running server..."
	@cd $(SERV_SRC_DIR) && $(GO) run -tags "$(SERVER_TAGS)" .

run-client:
	@echo "Running client..."
	@cd $(CLIENT_SRC_DIR) && $(GO) run .

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
