package main

import (
	"github.com/gin-gonic/gin"
	"github.com/telegram-go-stock-bot/pkg"
	"github.com/telegram-go-stock-bot/pkg/messages/discord"
	"github.com/telegram-go-stock-bot/pkg/messages/telegram"
)

func main() {
	pkg.InitEnv()
	telegram.Init()
	discord.Init()
	r := gin.Default()
	r.POST("/tg", telegram.Listener)
	r.POST("/dc/interactions", discord.Middleware, discord.Listener)
	r.Run()
}
