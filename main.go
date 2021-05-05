package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
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

	var updates tgbotapi.UpdatesChannel
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

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if strings.Contains(update.Message.Text, "/") {
			command := strings.Split(update.Message.Text, " ")[0][1:]
			if command == "" {
				return
			}
			args := strings.Split(update.Message.Text, " ")
			var arg string
			if len(args) > 1 {
				arg = args[1]
			}
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Println("Recovered in f", r)
					}
				}()
				re := regexp.MustCompile("[0-9]{4}$")
				if !re.MatchString(command) {
					set := &QueryResponse{}
					err := query(command, set)
					if err != nil {
						panic(err)
					}
					command = strings.ReplaceAll(set.ResultSet.Result[0].Symbol, ".TW", "")
				}

				var meta client.FugleAPIResponse
				var err error
				meta, err = myClient.Meta(command, false)
				if arg != "i" && err == nil {
					var quote client.FugleAPIResponse
					quote, err = myClient.Quote(command, false)
					if err != nil {
						panic(err)
					}
					meta.Data.Quote = quote.Data.Quote
				} else if err != nil {
					panic(err)
				}

				if err != nil || meta.Data.Info.SymbolID == "" || len(command) != 4 {
					msg.Text = "找不到此股票"
				} else if arg == "i" {
					msg.Text, err = convertByTemplate("meta", meta.Data)
					msg.ParseMode = "HTML"
				} else {
					msg.Text = convertQuote(meta.Data)
					msg.ParseMode = "MarkdownV2"
				}
				if err != nil {
					log.Panic(err)
				}
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}()
		}

	}
}
