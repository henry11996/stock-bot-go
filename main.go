package main

import (
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var Loc, _ = time.LoadLocation("Asia/Taipei")

func main() {
	bot, updates := boot()

	InitCache()
	InitSchedule()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s(%v)] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		c := make(chan string)

		args := strings.Split(update.Message.Text, " ")
		if len(args) > 1 {
			args = args[1:]
		} else {
			args = make([]string, 0)
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
	var Now = time.Now().In(Loc)

	fugle := fugleInit()

	defer func() {
		if r := recover(); r != nil {
			c <- "無法取得```" + command + "```"
			log.Println("Recovered in ", r)
		}
	}()

	stockId, _ := query(command)

	if len(args) > 0 && args[0] == "i" {
		meta, _ := fugle.Meta(stockId, false)
		c <- convertInfo(meta.Data)
	} else if len(args) > 0 && args[0] == "d" {
		releaseTime := time.Date(Now.Year(), Now.Month(), Now.Day(), 16, 2, 0, 0, Loc)
		if Now.Before(releaseTime) {
			Now = Now.AddDate(0, 0, -1)
		}
		var res LegalPersonResponse
		if len(args) > 1 {
			res = getDateLegalPersons(args[1])
		} else {
			res = getDateLegalPersons(Now.Format("20060102"))
		}
		legal := res.getByStock(stockId, command)
		c <- convertLegalPerson(legal)
	} else if len(args) > 0 && args[0] == "m" {
		var res LegalPersonResponse
		if len(args) > 1 {
			res = getMonthLegalPersons(args[1])
		} else {
			res = getMonthLegalPersons(Now.Format("20060102"))
		}
		legal := res.getByStock(stockId, command)
		c <- convertLegalPerson(legal)
	} else if command == "tw" {
		var legal LegalPerson
		if len(args) > 0 {
			legal = getDayTotalLegalPerson(args[0])
		} else {
			legal = getDayTotalLegalPerson(Now.Format("20060102"))
		}
		c <- convertTotalLegalPerson(legal)
	} else {
		meta, _ := fugle.Meta(stockId, false)
		quote, _ := fugle.Quote(stockId, false)
		meta.Data.Quote = quote.Data.Quote
		c <- convertQuote(meta.Data)
	}
}
