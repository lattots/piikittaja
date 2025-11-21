REMINDER_SRC_DIR := cmd/reminder
MESSENGER_SRC_DIR := cmd/messenger

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

