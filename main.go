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
		if len(args) > 1 {
			args = args[1:]
		}
		go run(strings.Split(update.Message.Text, " ")[0][1:], args, c)

		msg.Text = <-c
		msg.ParseMode = "MarkdownV2"
		log.Print(msg.Text)
		_, err := bot.Send(msg)
		if err != nil {
			log.Print(err)
		}
	}
}

func run(command string, args []string, c chan string) {
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
	if args[0] == "i" {
		c <- convertInfo(meta.Data)
	} else if args[0] == "d" {
		var res LegalPersonResponse
		if len(args) > 1 {
			res = getDateLegalPersons(args[1])
		} else {
			res = getDateLegalPersons("")
		}
		legal := res.getByStock(stockId, "")
		c <- convertLegalPerson(legal)
	} else if args[0] == "m" {
		var res LegalPersonResponse
		if len(args) > 1 {
			res = getMonthLegalPersons(args[1])
		} else {
			res = getMonthLegalPersons("")
		}
		legal := res.getByStock(stockId, "")
		c <- convertLegalPerson(legal)
	} else {
		quote, _ := fugle.Quote(stockId, false)
		meta.Data.Quote = quote.Data.Quote
		c <- convertQuote(meta.Data)
	}
}
