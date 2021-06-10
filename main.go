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

		args := strings.Split(update.Message.Text, " ")
		arg := ""
		if len(args) > 1 {
			arg = args[1]
		}
		go run(strings.Split(update.Message.Text, " ")[0][1:], arg, c)

		msg.Text = <-c
		msg.ParseMode = "MarkdownV2"
		log.Print(msg.Text)
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
	} else if arg == "t" {
		res := getLegalPersons("")
		legal := res.getByStock(stockId, "")
		log.Print(legal)
		c <- convertLegalPerson(legal)
	} else {
		quote, _ := fugle.Quote(stockId, false)
		meta.Data.Quote = quote.Data.Quote
		c <- convertQuote(meta.Data)
	}
}
