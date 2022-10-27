package botcode

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	ttlcache "github.com/jellydator/ttlcache/v3"
	"github.com/magiconair/properties"
)

const version string = "0.2"

var FIO string
var KEY string

var cacheBalance *ttlcache.Cache[string, Balance]

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

	cacheBalance = ttlcache.New(
		ttlcache.WithTTL[string, Balance](30 * time.Minute),
	)

	go cacheBalance.Start()

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
	case "help":
		answer.Text = help
	case "balance":
		answer.Text = getBalance(FIO)
	case "kill":
		fmt.Printf("Killed manually by %s \n", answer.ChannelUsername)
		os.Exit(0)
	default:
		answer.Text = "Unknown!"
	}
}

const help string = "ping, hi, ver"

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
