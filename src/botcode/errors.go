package botcode

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

var (
	master string
	bot    *tgbotapi.BotAPI
)

func InitErrorSender(masterId string, botApi *tgbotapi.BotAPI) {
	master = masterId
	bot = botApi
	go func() {
		_, err := bot.Send(tgbotapi.NewMessageToChannel(masterId, "hello"))
		if err != nil {
			fmt.Println(err)
		}
	}()
}

func SendError(error error) {
	if error != nil {
		SendErrorS(error.Error())
	} else {
		fmt.Println("error should be nil")
	}
}

func SendErrorS(error string) {
	if check() {
		_, err := bot.Send(tgbotapi.NewMessageToChannel(master, error))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func check() bool {
	if len(master) == 0 || bot == nil {
		log.Println("Cant send error to master. Init with 'InitErrorSender'")
		return false
	}
	return true
}
