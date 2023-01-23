package main

import (
	"log"
	"net/http"
	"os"

	"github.com/adshao/go-binance/v2"
	dsBot "github.com/bwmarrin/discordgo"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/henry11996/fugle-golang/fugle"
	"github.com/joho/godotenv"
)

func InitTelegramBot() (*tgBot.BotAPI, tgBot.UpdatesChannel) {
	bot, err := tgBot.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	// bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	var updates tgBot.UpdatesChannel
	if os.Getenv("LISTEN_MODE") == "socket" {
		_, _ = bot.SetWebhook(tgBot.NewWebhook(""))
		u := tgBot.NewUpdate(0)
		u.Timeout = 60
		updates, _ = bot.GetUpdatesChan(u)
	} else {
		url := os.Getenv("TELEGRAM_WEBHOOK_URL") + "/" + bot.Token
		_, err = bot.SetWebhook(tgBot.NewWebhook(url))
		if err != nil {
			log.Fatal(err)
		}
		var info tgBot.WebhookInfo
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
	return bot, updates
}

func InitFugle() fugle.Client {
	client, err := fugle.NewFugleClient(fugle.ClientOption{
		ApiToken: os.Getenv("FUGLE_API_TOKEN"),
		Version:  "v0.3",
	})
	if err != nil {
		log.Fatal("failed to init fugle api client, " + err.Error())
	}
	return client
}

func InitBinance() *binance.Client {
	var (
		apiKey    = os.Getenv("BINANCE_API_KEY")
		secretKey = os.Getenv("BINANCE_SECRET_KEY")
	)
	return binance.NewClient(apiKey, secretKey)
}

func InitEnv() error {
	return godotenv.Load()
}

func InitDiscord() *dsBot.Session {
	discord, err := dsBot.New("Bot " + "authentication token")
	if err != nil {
		log.Fatal("failed to init discord api client, " + err.Error())
	}
	return discord
}
