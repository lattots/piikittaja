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
all: $(MANAGER_BIN)

$(MANAGER_BIN): $(wildcard $(MANAGER_SRC_DIR)/*.go) $(wildcard pkg/**/*.go)
	mkdir -p $(BIN_DIR)
	go build -o $@ $(MANAGER_SRC_DIR)/manager.go

# Build target for the Github build action
build: all

# Runs docker compose that spins up containers with "piikki-web" and "piikki-bot" images from Docker Hub
# Images need to be pushed to the repository before running this
compose-up:
	@docker compose -f ./cicd/compose.yaml up -d

compose-down: clean-bot clean-web

# ---
# Web app commands

.PHONY: run-web stop-web clean-web log-web deploy-web

build-web: $(wildcard $(TELEGRAM_BOT_SRC_DIR)/*.go) $(wildcard pkg/**/*.go)
	@docker build -t lattots/piikki-web -f ./cicd/web_app/Dockerfile .

run-web: build-web
	@docker run -d --network="host" --name web-app-container lattots/piikki-web

stop-web:
	@docker stop web-app-container

clean-web: stop-web
	@docker rm web-app-container

log-web:
	@docker logs web-app-container

deploy-web: build-web
	@docker push lattots/piikki-web:latest

# ---
# Telegram bot commands

.PHONY: run-bot stop-bot clean-bot log-bot deploy-bot

build-bot: $(wildcard $(TELEGRAM_BOT_SRC_DIR)/*.go) $(wildcard pkg/**/*.go)
	@docker build -t lattots/piikki-bot -f ./cicd/telegram_bot/Dockerfile .

run-bot: build-bot
	@docker run -d --network="host" --name telegram-bot-container lattots/piikki-bot

stop-bot:
	@docker stop telegram-bot-container

clean-bot: stop-bot
	@docker rm telegram-bot-container

log-bot:
	@docker logs telegram-bot-container

deploy-bot: build-bot
	@docker push lattots/piikki-bot:latest

# ---
# Admin manager commands

.PHONY: manager
manager: $(MANAGER_BIN)
	$(MANAGER_BIN)

# ---
# Payment reminder commands

.PHONY: remind run-reminder clean-reminder log-reminder deploy-reminder

remind: stop-bot run-reminder compose-up

build-reminder: $(wildcard $(REMINDER_SRC_DIR)/*.go) $(wildcard pkg/**/*.go)
	@docker build -t lattots/piikki-reminder -f ./cicd/reminder/Dockerfile .

<<<<<<< Updated upstream
run-reminder: build-reminder
	@docker run -d --network="host" --name reminder-container lattots/piikki-reminder
=======
run-reminder:
	@docker run -d --pull always --network="host" --name reminder-container lattots/piikki-reminder
>>>>>>> Stashed changes

clean-reminder:
	@docker rm reminder-container

log-reminder:
	@docker logs reminder-container

deploy-reminder: build-reminder
	@docker push lattots/piikki-reminder:latest

# Clean up binaries
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)/*
