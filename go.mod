module github.com/and-n/telegBot

go 1.12

replace github.com/and-n/telegBot/botcode v1.0.0 => ./src/botcode

require (
	github.com/and-n/telegBot/botcode v1.0.0
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
)
