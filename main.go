package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/RainrainWu/fugle-realtime-go/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析參數，預設是不會解析的
}

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

	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60
	// updates, err := bot.GetUpdatesChan(u)

	url := os.Getenv("HEROKU_APP_NAME") + ".herokuapp.com" + "/"
	log.Printf(url)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(url + bot.Token))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "HTML"
		if update.Message.IsCommand() {
			command := update.Message.Command()
			if command == "" {
				return
			}
			args := update.Message.CommandArguments()

			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Println("Recovered in f", r)
					}
				}()

				re := regexp.MustCompile("[0-9]{4}$")
				if !re.MatchString(command) {
					return
				}

				var data client.FugleAPIResponse
				var err error
				if args == "i" {
					data, err = myClient.Meta(command, false)
				} else {
					data, err = myClient.Quote(command, false)
				}

				if err != nil || data.Data.Info.SymbolID == "" || len(command) != 4 {
					msg.Text = "找不到此股票"
				} else if args == "i" {
					msg.Text, err = convertByTemplate("meta", data.Data)
				} else {
					msg.Text, err = convertByTemplate("quote", data.Data)
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
