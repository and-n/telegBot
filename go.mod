module github.com/and-n/telegBot

go 1.12

replace github.com/userName/otherModule v0.0.0 => ./

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 // indirect
	github.com/magiconair/properties v1.8.6
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	golang.org/x/text v0.4.0
)
