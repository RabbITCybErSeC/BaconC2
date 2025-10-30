GO = go
BINARY_NAME = bacon-server
CLIENT_BINARY_NAME = bacon-client
SERV_SRC_DIR = ./server
CLIENT_SRC_DIR = ./client
FRONT_END_DIR = ./becongui
BIN_DIR = ./bin
BUILD_CONFIG = build.json

PLATFORMS = darwin linux windows
ARCHS = amd64 arm64
TARGET_OS ?= darwin
TARGET_ARCH ?= arm64
TARGET = $(TARGET_OS)-$(TARGET_ARCH)

GET_TAGS = $(shell jq -r '.components[] | select(.enabled) | .build_tag' $(BUILD_CONFIG) | tr '\n' ' ')

VALID_OS := $(filter $(TARGET_OS),$(PLATFORMS))
VALID_ARCH := $(filter $(TARGET_ARCH),$(ARCHS))
ifneq ($(VALID_OS),$(TARGET_OS))
  $(error Invalid TARGET_OS: $(TARGET_OS). Must be one of: $(PLATFORMS))
endif
ifneq ($(VALID_ARCH),$(TARGET_ARCH))
  $(error Invalid TARGET_ARCH: $(TARGET_ARCH). Must be one of: $(ARCHS))
endif

# Build tags and flags
LDFLAGS = -s -w
SERVER_TAGS = go_sqlite,server,$(TARGET_OS) $(GET_TAGS)
CLIENT_TAGS = $(TARGET_OS) $(GET_TAGS)

SERVER_BUILD_FLAGS = -trimpath -tags "$(SERVER_TAGS)" -ldflags "$(LDFLAGS)"
CLIENT_BUILD_FLAGS = -trimpath -tags "$(CLIENT_TAGS)" -ldflags "$(LDFLAGS)"

SERVER_OUTPUT = $(BIN_DIR)/$(BINARY_NAME)_$(TARGET)
CLIENT_OUTPUT = $(BIN_DIR)/$(CLIENT_BINARY_NAME)_$(TARGET)

.PHONY: all all-platforms clean compile-server compile-client run-server run-client build-front-end show-modules

all: compile-server compile-client

all-platforms: build-front-end
	@echo "Building for all platforms..."
	@for os in $(PLATFORMS); do \
		for arch in $(ARCHS); do \
			$(MAKE) all TARGET_OS=$$os TARGET_ARCH=$$arch || exit 1; \
		done; \
	done
	@echo "All platforms built successfully"

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

build-front-end:
	@echo "Building front-end"
	@cd $(FRONT_END_DIR) && npm run build

compile-server: $(BIN_DIR)
	@echo "Building server for $(TARGET)..."
	@echo "→ Tags: $(SERVER_TAGS)"
	@cd $(SERV_SRC_DIR) && GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) CGO_ENABLED=1 $(GO) build $(SERVER_BUILD_FLAGS) -o ../$(SERVER_OUTPUT) .
	@echo "Server build completed: $(SERVER_OUTPUT)"

compile-client: $(BIN_DIR)
	@echo "Building client for $(TARGET)..."
	@echo "→ Tags: $(CLIENT_TAGS)"
	@cd $(CLIENT_SRC_DIR) && GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) CGO_ENABLED=0 $(GO) build $(CLIENT_BUILD_FLAGS) -o ../$(CLIENT_OUTPUT) .
	@echo "Client build completed: $(CLIENT_OUTPUT)"

run-server: build-front-end
	@echo "Running server..."
	@cd $(SERV_SRC_DIR) && $(GO) run -tags "$(SERVER_TAGS)" .

run-client:
	@echo "Running client..."
	@cd $(CLIENT_SRC_DIR) && $(GO) run -tags "$(CLIENT_TAGS)" .

show-modules:
	@echo "\nEnabled components from build.json:"
	@jq -r '.components[] | select(.enabled) | "- " + .name + " (tag: " + .build_tag + ")"' $(BUILD_CONFIG)
	@echo

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)