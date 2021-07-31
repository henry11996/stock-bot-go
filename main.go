package main

import (
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var Loc, _ = time.LoadLocation("Asia/Taipei")

func main() {
	bot, updates := Boot()

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
			err := ""
			switch x := r.(type) {
			case string:
				err = x
			case error:
				err = x.Error()
			}
			c <- "無法取得```" + command + "```，" + err
		}
	}()

	stockId, _ := query(command)
	var err error
	if command == "tw" {
		if len(args) > 0 && args[0] == "d" {
			t := Now
			if len(args) > 1 {
				t, err = time.Parse("2006/01/02", args[1])
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm/dd")
				}
			}
			lp, err := getDayTotalLegalPerson(t)
			if err != nil {
				log.Panic(err)
			}
			c <- lp.PrettyString()
		} else if len(args) > 0 && args[0] == "m" {
			t := Now
			if len(args) > 1 {
				t, err = time.Parse("2006/01", args[1])
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm")
				}
			}
			lp, err := getMonthTotalLegalPerson(t)
			if err != nil {
				log.Panic(err)
			}
			c <- lp.PrettyString()
		} else {
			log.Panic("錯誤指令")
		}
	} else {
		if len(args) > 0 && args[0] == "i" {
			meta, _ := fugle.Meta(stockId, false)
			c <- convertInfo(meta.Data)
		} else if len(args) > 0 && args[0] == "d" {
			t := Now
			if len(args) > 1 {
				t, err = time.Parse("2006/01/02", args[1])
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm/dd")
				}
			}
			lp, _ := getDayLegalPersons(t)

			c <- lp.FindStock(stockId, command).PrettyString(lp.Title)
		} else if len(args) > 0 && args[0] == "m" {
			t := Now
			if len(args) > 1 {
				t, err = time.Parse("2006/01", args[1])
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm")
				}
			}
			lp, _ := getMonthLegalPersons(t)

			c <- lp.FindStock(stockId, command).PrettyString(lp.Title)
		} else {
			meta, _ := fugle.Meta(stockId, false)
			quote, _ := fugle.Quote(stockId, false)
			meta.Data.Quote = quote.Data.Quote
			c <- convertQuote(meta.Data)
		}
	}
}
