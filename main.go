package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, updates := boot()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		c := make(chan string)

		go run(strings.Split(update.Message.Text, " ")[0][1:], update.Message.CommandArguments(), c)

		msg.Text = <-c
		msg.ParseMode = "MarkdownV2"
		_, err := bot.Send(msg)
		if err != nil {
			log.Print(err)
		}
	}
}

func run(command string, arg string, c chan string) {
	fugle := fugleInit()

	defer func() {
		if r := recover(); r != nil {
			c <- "找不到```" + command + "```"
			log.Println("Recovered in ", r)
		}
	}()
	if command == "twse" {
		command = "TWSE_SEM_INDEX_1"
	}

	stockId, _ := query(command)

	meta, _ := fugle.Meta(stockId, false)
	if arg == "i" {
		c <- convertInfo(meta.Data)
	} else {
		quote, _ := fugle.Quote(stockId, false)
		meta.Data.Quote = quote.Data.Quote
		c <- convertQuote(meta.Data)
	}
}
