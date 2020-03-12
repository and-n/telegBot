package main

import (
	"flag"
	"fmt"
	"github.com/and-n/telegBot/botcode"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"os/signal"
)

func main() {
	fmt.Println("hello!")
	var apiKey, masterId string
	flag.StringVar(&apiKey, "key", "", "api key for telega")
	flag.StringVar(&masterId, "master", "", "master id")
	flag.Parse()

	bot, updates := botcode.InitBot(apiKey)
	botcode.InitErrorSender(masterId, bot)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		botcode.SendErrorS("normal Shutting down")
		bot.StopReceivingUpdates()
		os.Exit(0)
	}()

	defer botcode.SendErrorS("Failed Shutting down")
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		go botcode.SaveStatistic(update.Message.From)

		go func(api *tgbotapi.BotAPI) {
			msg := botcode.AnswerMessage(update.Message)
			_, err := bot.Send(msg)
			if err != nil {
				log.Fatal(err)
			}
		}(bot)
	}
}
