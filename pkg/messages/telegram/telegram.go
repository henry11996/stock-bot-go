package telegram

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telegram-go-stock-bot/pkg"
)

var Bot *tgBot.BotAPI

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Init() *tgBot.BotAPI {
	var token = os.Getenv("TELEGRAM_API_TOKEN")
	var err error
	Bot, err = tgBot.NewBotAPI(token)
	Check(err)

	url := os.Getenv("TELEGRAM_WEBHOOK_URL") + "/tg" //"/" + bot.Token +
	_, err = Bot.SetWebhook(tgBot.NewWebhook(url))
	Check(err)

	log.Printf("Authorized on account %s", Bot.Self.UserName)
	return Bot
}

func Listener(c *gin.Context) {
	bytes, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	var update tgBot.Update
	json.Unmarshal(bytes, &update)

	if update.Message == nil {
		return
	}
	log.Printf("[%s(%v)] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

	command := strings.Split(update.Message.Text, " ")[0][1:]
	args := strings.Split(update.Message.Text, " ")
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = make([]string, 0)
	}

	ch := make(chan string)
	go pkg.Route(command, args, ch)
	message := <-ch

	msg := tgBot.NewMessage(update.Message.Chat.ID, "")
	msg.Text = message
	msg.ParseMode = "MarkdownV2"
	_, err := Bot.Send(msg)
	if err != nil {
		log.Print(err)
	}
}
