package main

import (
	"log"
	"os"
	"strings"

	"github.com/RainrainWu/fugle-realtime-go/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	myClient, err := client.NewFugleClient()
	if err != nil {
		log.Fatal("failed to init fugle api client")
	}

	u := tgbotapi.NewUpdate(1)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "HTML"

		if strings.Contains(update.Message.Text, "/") {
			go func() {
				defer func() {
					if err := recover(); err != nil {
						log.Println("recover success. " + update.Message.Text)
					}
				}()
				i := strings.Index(update.Message.Text, "/")
				result := myClient.Meta(update.Message.Text[i+1:i+5], false)
				if result.Data.Info.SymbolID == "" {
					msg.Text = "找不到此股票"
				} else {
					msg.Text = GetInfo(result.Data)
				}
				bot.Send(msg)
			}()
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		arg := update.Message.CommandArguments()

		switch update.Message.Command() {
		case "tw":
			text := ""
			go func() {
				defer func() {
					if err := recover(); err != nil {
						log.Println("recover success.")
					}
				}()
				if arg == "" || len(arg) != 4 {
					text = "北七？"
				} else {
					result := myClient.Meta(arg, false)
					text = GetInfo(result.Data)
				}
			}()
			msg.Text = text
		case "us":
		case "status":
			msg.Text = "I'm ok."
		default:
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if _, err := bot.Send(msg); err != nil {
			failMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error occur "+err.Error())
			bot.Send(failMsg)
		}
	}
}
