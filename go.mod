module github.com/and-n/telegBot

go 1.12

replace github.com/and-n/telegBot/botcode v1.0.0 => ./src/botcode

require (
	github.com/and-n/telegBot/botcode v1.0.0
	gopkg.in/yaml.v3 v3.0.0-20200121175148-a6ecf24a6d71
)
