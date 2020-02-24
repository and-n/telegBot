package main

import (
	"fmt"
	"github.com/and-n/telegBot/botcode"
	"log"
	"os"
)

func main() {
	fmt.Println("hello!")

	apiKey := os.Getenv("API_KEY")
	if len(apiKey) == 0 {
		apiKey = os.Args[1]
		if len(apiKey) == 0 {
			log.Fatal("dont forget to set APIkey")
		}
	}

	bot, updates := botcode.InitBot(apiKey)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		go botcode.SaveStatistic(update.Message.From)

		msg := botcode.AnswerMessage(update.Message)

		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}

}
