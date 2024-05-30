run:
	./bin/telegram_bot

build:
	go build -o bin/telegram_bot src/bot_server/server.go
