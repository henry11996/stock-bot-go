package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stock-bot-go/pkg"
	"github.com/stock-bot-go/pkg/messages/discord"
	"github.com/stock-bot-go/pkg/messages/telegram"
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
