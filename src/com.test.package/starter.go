package main

import (
	"com.test.package/botcode"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("hello!")
	fmt.Println("dont forget to set APIkey")

	apiKey := os.Args[1]

	bot, updates := botcode.InitBot(apiKey)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := botcode.AnswerMessage(update.Message)

		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}

}