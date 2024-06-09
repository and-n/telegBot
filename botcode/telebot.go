package botcode

import (
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/magiconair/properties"
)

const version string = "0.4"

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

func AnswerInlineQuery(query *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	println(query.Data)
	splitted := strings.Split(query.Data, ":")
	if len(splitted) > 1 {
		if splitted[0] == "month" {
			month, _ := strconv.ParseInt(splitted[1], 10, 0)
			sum, err := getSumByMonthAsString(FIO, int(month))
			if err != nil {
				bot.Send(tgbotapi.NewMessage(query.From.ID, err.Error()))
			} else {
				bot.Send(tgbotapi.NewMessage(query.From.ID, time.Month(month).String()+":\n"+sum))
			}
		} else {
			bot.Send(tgbotapi.NewMessage(query.From.ID, "error"))
		}
	}
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
	case "month":
		if len(arguments) == 0 {
			answer.Text = "/month MONTH_NUMBER"
		} else {
			m, err := strconv.ParseInt(arguments, 10, 0)
			if err != nil {
				answer.Text = "Wrong month number"
			} else {
				res, err := getSumByMonthAsString(FIO, int(m))
				if err != nil {
					answer.Text = err.Error()
				} else {
					answer.Text = res
				}
			}

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
	case "month":
		answer.Text = "/month MONTH_NUMBER"
		answer.ReplyMarkup = getMonthKeyboard()
	default:
		answer.Text = "Unknown!"
	}

}

func createButtons() tgbotapi.ReplyKeyboardMarkup {
	buttons := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("balance"),
			tgbotapi.NewKeyboardButton("month"),
			tgbotapi.NewKeyboardButton("help"),
		),
	)
	return buttons
}

func getMonthKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(time.Now().Month().String(), "month:"+strconv.Itoa(int(time.Now().Month()))),
			tgbotapi.NewInlineKeyboardButtonData(monthChange(time.Now().Month(), -1).String(), "month:"+strconv.Itoa(int(monthChange(time.Now().Month(), -1)))),
			tgbotapi.NewInlineKeyboardButtonData(monthChange(time.Now().Month(), -2).String(), "month:"+strconv.Itoa(int(monthChange(time.Now().Month(), -2)))),
		),
	)
	return keyboard
}

func monthChange(month time.Month, change int) time.Month {
	if change == 0 {
		return month
	}
	var newMonth int
	if change > 0 {
		newMonth = int(month) + change%12
		if int(newMonth) <= 12 {
			return time.Month(newMonth)
		} else {
			return time.Month(newMonth % 12)
		}
	} else {
		newMonth = int(month) + change%12
		if newMonth > 0 {
			return time.Month(newMonth)
		} else {
			return time.Month(12 + newMonth)
		}
	}

}
