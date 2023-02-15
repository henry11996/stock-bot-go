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
	r.POST("/_ah/warmup", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"ok": true,
		})
	})
	r.POST("/tg", telegram.Listener)
	r.POST("/dc/interactions", discord.Middleware, discord.Listener)
	r.Run()
}
