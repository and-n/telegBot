package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/and-n/telegBot/botcode"
	"github.com/magiconair/properties"
)

func main() {
	fmt.Println("hello!")

	if len(os.Args) > 2 {
		log.Fatal("strange args! it is possible to setup properties file and that it! %n", len(os.Args))
	}
	var p *properties.Properties
	if len(os.Args) == 2 {
		p = properties.MustLoadFile(os.Args[1], properties.UTF8)
	} else {
		p = properties.MustLoadFile("./configuration/file.properties", properties.UTF8)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		if sig == os.Interrupt {
			os.Exit(0)
		}
	}()

	bot, updates := botcode.InitBot(p)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		go botcode.SaveStatistic(update.Message.From)

		go botcode.AnswerMessage(update.Message, bot)

	}

}
