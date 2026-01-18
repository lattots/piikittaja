REMINDER_SRC_DIR := cmd/reminder
MESSENGER_SRC_DIR := cmd/messenger

up:
	docker compose -f compose.yaml up --build

upd:
	docker compose -f compose.yaml up --build -d

down:
	docker compose down

dev:
	docker compose -f dev.compose.yaml up -d
	docker compose -f dev.compose.yaml watch

ref:
	docker compose -f dev.compose.yaml up -d --build web_builder web_app

dev_down:
	docker compose -f ./dev.compose.yaml down --volumes --remove-orphans

test:
	docker compose -f ./test.compose.yaml down --volumes --remove-orphans
	docker compose -f ./test.compose.yaml up --build

