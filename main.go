package main

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	tgBot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/freetype/truetype"
)

var Loc, _ = time.LoadLocation("Asia/Taipei")
var DefaultFont = readFont()

func main() {
	InitEnv()
	InitCache()
	bot, updates := InitTelegramBot()
	listenTelegramUpdateChannel(bot, updates)
}

func listenTelegramUpdateChannel(bot *tgBot.BotAPI, updates tgBot.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s(%v)] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

		c := make(chan interface{})

		command := strings.Split(update.Message.Text, " ")[0][1:]
		args := strings.Split(update.Message.Text, " ")
		if len(args) > 1 {
			args = args[1:]
		} else {
			args = make([]string, 0)
		}

		go Route(command, args, c)

		message := <-c
		switch v := message.(type) {
		case string:
			msg := tgBot.NewMessage(update.Message.Chat.ID, "")
			msg.Text = v
			msg.ParseMode = "MarkdownV2"
			_, err := bot.Send(msg)
			if err != nil {
				log.Print(err)
			}
		case []byte:
			msg := tgBot.NewPhotoUpload(update.Message.Chat.ID, tgBot.FileBytes{
				Bytes: v,
			})
			_, err := bot.Send(msg)
			if err != nil {
				log.Print(err)
			}
		default:
		}
	}
}

func readFont() *truetype.Font {
	b, err := ioutil.ReadFile("./fonts/TaipeiSansTCBeta-Bold.ttf")
	if err != nil {
		log.Panic(err)
	}
	font, err := truetype.Parse(b)
	if err != nil {
		log.Panic(err)
	}
	return font
}
