run:
	./bin/telegram_bot

build:
	go build -o bin/telegram_bot cmd/telegram_bot/server.go
	go build -o bin/web_app cmd/web_app/server.go
