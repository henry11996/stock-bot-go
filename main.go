package main

import (
	"log"
	"os"
	"strconv"

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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		twStockId := ""
		// usStockId := ""
		switch update.Message.Command() {
		case "tw":
			twStockId = update.Message.CommandArguments()
			result := myClient.Meta(twStockId, false)
			result.PrettyPrint()
			msg.Text = strconv.Itoa(result.Data.Meta.Pricereference)
		case "us":
			// usStockId = update.Message.CommandArguments()
		case "status":
			msg.Text = "I'm ok."
		default:
			msg.Text = "I don't know that command"
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}
