package botcode

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/magiconair/properties"
)

const version string = "0.3"

var FIO string
var KEY string

// InitBot -init telegram bot
func InitBot(props *properties.Properties) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	loadProperties(props)

	bot, err := tgbotapi.NewBotAPI(KEY)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	//return bot
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	return bot, updates
}

func loadProperties(props *properties.Properties) {
	KEY = props.MustGetString("bot_api")
	if len(KEY) == 0 {
		log.Panic("API KEY!")
	}
	FIO = props.MustGetString("fio_api")
	if len(FIO) == 0 {
		log.Panic("FIO KEY!")
	}
}

func AnswerMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	var answer tgbotapi.MessageConfig
	answer.ChatID = message.Chat.ID
	answer.ChannelUsername = message.From.UserName
	answer.Text = "kill me please"

	answer.ReplyMarkup = createButtons()

	if message.IsCommand() {
		parseCommand(message.Command(), message.CommandArguments(), &answer)
	} else if len(message.Text) != 0 {
		answer.ReplyToMessageID = message.MessageID
		parseString(message, &answer)
	}

	bot.Send(answer)
}

func parseCommand(command string, arguments string, answer *tgbotapi.MessageConfig) {
	// log.Println("command", command, "arg ", arguments, "at", at)
	switch command {

	case "balance":
		answer.Text = getBalance(FIO)
	case "api":
		if len(arguments) == 0 {
			answer.Text = "/api YOUR_KEY"
		} else {
			answer.Text = "your api: " + arguments
		}

	// case "kill":
	// 	fmt.Printf("Killed manually by %s \n", answer.ChannelUsername)
	// 	os.Exit(0)
	default:
		answer.Text = "Unknown!"
	}
}

const help string = "ping, hi, ver, balance"

func parseString(message *tgbotapi.Message, answer *tgbotapi.MessageConfig) {

	switch strings.ToLower(message.Text) {
	case "ping":
		answer.Text = "pong"
	case "hello", "hi":
		answer.Text = "Hello, " + message.From.UserName
	case "ver", "version":
		answer.Text = version
	case "balance":
		answer.Text = getBalance(FIO)

	case "api":
		answer.Text = "/api YOUR_KEY"
		answer.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	case "help":
		answer.Text = help
	default:
		answer.Text = "Unknown!"
	}

}

func createButtons() tgbotapi.ReplyKeyboardMarkup {
	buttons := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("balance"),
			tgbotapi.NewKeyboardButton("help"),
		),
	)
	return buttons
}
