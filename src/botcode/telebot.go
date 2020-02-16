package botcode

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

const version string = "0.1"

// InitBot -init telegram bot
func InitBot(key string) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	if len(key) == 0 {
		log.Panic("API KEY!")
	}
	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	//return bot
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	return bot, updates
}

func AnswerMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	var answer tgbotapi.MessageConfig
	answer.ChatID = message.Chat.ID
	answer.ReplyToMessageID = message.MessageID
	answer.Text = "kill me please"

	if message.IsCommand() {
		parseCommand(message.Command(), message.CommandArguments(), message.CommandWithAt(), &answer)
	} else if "" != message.Text {
		parseString(message, &answer)
	}
	return answer
}

func parseCommand(command string, arguments string, at string, answer *tgbotapi.MessageConfig) {
	log.Println("command", command, "arg ", arguments, "at", at)
}

func parseString(message *tgbotapi.Message, answer *tgbotapi.MessageConfig) {
	switch strings.ToLower(message.Text) {
	case "ping":
		answer.Text = "pong"
	case "hello", "hi":
		answer.Text = "Hello, " + message.From.UserName
	case "ver", "version":
		answer.Text = version
	default:
		answer.Text = "Sorry"
	}
}
