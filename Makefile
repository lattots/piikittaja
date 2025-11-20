# Source directories (assuming all related Go files are under these directories)
TELEGRAM_BOT_SRC_DIR := cmd/telegram_bot
WEB_APP_SRC_DIR := cmd/web_app
REMINDER_SRC_DIR := cmd/reminder
MANAGER_SRC_DIR := cmd/admin_manager
MESSENGER_SRC_DIR := cmd/messenger

# Runs docker compose that spins up containers with "piikki-web" and "piikki-bot" images from Docker Hub
# Images need to be pushed to the repository before running this
up:
	docker compose -f compose.yaml up --build

upd:
	docker compose -f compose.yaml up --build -d

down:
	docker compose down

test:
	docker compose -f ./test.compose.yaml down --volumes --remove-orphans
	docker compose -f ./test.compose.yaml up --build

# ---
# Admin manager commands

.PHONY: run-manager clean-manager log-manager deploy-manager

build-manager: $(wildcard $(MANAGER_SRC_DIR)/*.go) $(wildcard pkg/**/*.go)
	@docker build -t lattots/piikki-manager -f ./cicd/admin_manager/Dockerfile .

run-manager:
	@docker run -it --rm --pull always --network="host" --name manager-container lattots/piikki-manager

clean-manager:
	@docker rm manager-container

log-manager:
	@docker logs manager-container

deploy-manager: build-manager
	@docker push lattots/piikki-manager:latest

# ---
# Payment reminder commands

.PHONY: remind run-reminder clean-reminder log-reminder deploy-reminder

remind: stop-bot run-reminder compose-up

build-reminder: $(wildcard $(REMINDER_SRC_DIR)/*.go) $(wildcard pkg/**/*.go)
	@docker build -t lattots/piikki-reminder -f ./cicd/reminder/Dockerfile .

run-reminder:
	@docker run -d --pull always --network="host" --name reminder-container lattots/piikki-reminder

clean-reminder:
	@docker rm reminder-container

log-reminder:
	@docker logs reminder-container

deploy-reminder: build-reminder
	@docker push lattots/piikki-reminder:latest

# ---
# User messenger commands

.PHONY: messenger run-messenger clean-messenger log-messenger deploy-messenger

messenger: stop-bot run-messenger compose-up

build-messenger: $(wildcard $(MESSENGER_SRC_DIR)/*.go) $(wildcard pkg/**/*.go)
	@docker build -t lattots/piikki-messenger -f ./cicd/messenger/Dockerfile .

run-messenger:
	@docker run -d --pull always --network="host" --name messenger-container lattots/piikki-messenger

clean-messenger:
	@docker rm messenger-container

log-messenger:
	@docker logs messenger-container

deploy-messenger: build-messenger
	@docker push lattots/piikki-messenger:latest

