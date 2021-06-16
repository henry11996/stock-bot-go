package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron"
)

func InitSchedule() {
	c := cron.New()
	c.AddFunc("0 3 7 * * 1-5", setDailyLegalPerson)
	c.Start()
}

func setDailyLegalPerson() {
	log.Print("開始快取當日所有股票三大法人買日買賣超")
	legal := getDayTotalLegalPerson("")
	text := convertTotalLegalPerson(legal)
	groupId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_GROUP_ID"), 10, 64)
	msg := tgbotapi.NewMessage(groupId, text)
	msg.ParseMode = "MarkdownV2"
	_, err := Bot.Send(msg)
	if err != nil {
		log.Print(err)
	}
}
