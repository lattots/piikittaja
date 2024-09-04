# Directories and binaries
BIN_DIR := bin
TELEGRAM_BOT_BIN := $(BIN_DIR)/telegram_bot
WEB_APP_BIN := $(BIN_DIR)/web_app
REMINDER_BIN := $(BIN_DIR)/reminder
MANAGER_BIN := $(BIN_DIR)/manager

# Source directories (assuming all related Go files are under these directories)
TELEGRAM_BOT_SRC_DIR := cmd/telegram_bot
WEB_APP_SRC_DIR := cmd/web_app
REMINDER_SRC_DIR := cmd/reminder
MANAGER_SRC_DIR := cmd/admin_manager

# Default target
all: $(TELEGRAM_BOT_BIN) $(WEB_APP_BIN) $(REMINDER_BIN) $(MANAGER_BIN)

# Build targets
$(TELEGRAM_BOT_BIN): $(wildcard $(TELEGRAM_BOT_SRC_DIR)/*.go) $(wildcard pkg/**/*.go)
	mkdir -p $(BIN_DIR)
	go build -o $@ $(TELEGRAM_BOT_SRC_DIR)/server.go

$(WEB_APP_BIN): $(wildcard $(WEB_APP_SRC_DIR)/*.go) $(wildcard pkg/**/*.go)
	mkdir -p $(BIN_DIR)
	go build -o $@ $(WEB_APP_SRC_DIR)/server.go

$(REMINDER_BIN): $(wildcard $(REMINDER_SRC_DIR)/*.go) $(wildcard pkg/**/*.go)
	mkdir -p $(BIN_DIR)
	go build -o $@ $(REMINDER_SRC_DIR)/reminder.go

$(MANAGER_BIN): $(wildcard $(MANAGER_SRC_DIR)/*.go) $(wildcard pkg/**/*.go)
	mkdir -p $(BIN_DIR)
	go build -o $@ $(MANAGER_SRC_DIR)/manager.go

# Run commands
.PHONY: run
run: all
	$(TELEGRAM_BOT_BIN) &
	$(WEB_APP_BIN) &

.PHONY: run-web-app
run-web-app: $(WEB_APP_BIN)
	$(WEB_APP_BIN)

.PHONY: run-telegram-bot
run-telegram-bot: $(TELEGRAM_BOT_BIN)
	$(TELEGRAM_BOT_BIN)

.PHONY: manager
manager: $(MANAGER_BIN)
	$(MANAGER_BIN)

.PHONY: remind
remind: $(REMINDER_BIN)
	$(REMINDER_BIN)

# Clean up binaries
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)/*

# Production run
.PHONY: prod
prod: all
	$(TELEGRAM_BOT_BIN) &
	$(WEB_APP_BIN) &
