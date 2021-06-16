package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RainrainWu/fugle-realtime-go/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var Bot *tgbotapi.BotAPI

func boot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	var err error
	Bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}
	// bot.Debug = true
	log.Printf("Authorized on account %s", Bot.Self.UserName)

	return Bot, botInit(Bot)
}

func botInit(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	var updates tgbotapi.UpdatesChannel
	var err error
	if os.Getenv("APP_ENV") == "debug" {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(""))
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates, err = bot.GetUpdatesChan(u)
	} else {
		url := os.Getenv("HEROKU_APP_NAME") + ".herokuapp.com" + "/"
		log.Printf(url)
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(url + bot.Token))
		if err != nil {
			log.Fatal(err)
		}
		var info tgbotapi.WebhookInfo
		info, err = bot.GetWebhookInfo()
		if err != nil {
			log.Fatal(err)
		}
		if info.LastErrorDate != 0 {
			log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
		}
		updates = bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), nil)
		log.Printf("Server up with port " + os.Getenv("PORT"))
	}
	return updates
}

func fugleInit() client.FugleClient {
	myClient, err := client.NewFugleClient()
	if err != nil {
		log.Fatal("failed to init fugle api client")
	}
	return myClient
}
