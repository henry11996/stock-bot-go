package discord

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/telegram-go-stock-bot/pkg"
)

var (
	_ = pkg.InitEnv()
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Init() *discordgo.Session {
	var (
		token     = os.Getenv("DISCORD_API_TOKEN")
		appId     = os.Getenv("DISCORD_APP_ID")
		publicKey = os.Getenv("DISCORD_PUBLIC_KEsY")
	)
	// fmt.Print(publicKey)
	discord, err := discordgo.New("Bot " + token)
	fmt.Println(token, appId, publicKey)
	check(err)

	err = removeAllCommands(discord, appId)
	check(err)

	err = registerCommand(discord, appId)
	check(err)

	log.Printf("Discord command registered")

	return discord
}

func Listener(c *gin.Context) {
	var request discordgo.Interaction
	err := c.BindJSON(&request)
	check(err)
	if request.Type.String() == discordgo.InteractionPing.String() {
		c.JSON(http.StatusOK, gin.H{
			"type": 1,
		})
		return
	}

	var (
		command string
		args    []string
	)
	for _, option := range request.ApplicationCommandData().Options {
		if option.Name == "stock" {
			command = fmt.Sprintf("%v", option.Value)
		} else {
			args = append(args, fmt.Sprintf("%v", option.Value))
		}
	}
	ch := make(chan string)
	go pkg.Route(command, args, ch)
	message := <-ch
	message = strings.ReplaceAll(message, "\\", "")

	c.JSON(http.StatusOK, discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: message},
	})
}

func Middleware(c *gin.Context) {
	publicKey := os.Getenv("DISCORD_PUBLIC_KEY")
	publicKeyBytes, _ := hex.DecodeString(publicKey)

	verify := discordgo.VerifyInteraction(c.Request, publicKeyBytes)
	if !verify {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Next()
}

func removeAllCommands(discord *discordgo.Session, appId string) (err error) {
	commands, err := discord.ApplicationCommands(appId, "")
	for _, command := range commands {
		discord.ApplicationCommandDelete(command.ApplicationID, "", command.ID)
	}
	return
}

func registerCommand(discord *discordgo.Session, appId string) (err error) {
	_, err = discord.ApplicationCommandCreate(appId, "", &discordgo.ApplicationCommand{
		Name:        "tw",
		Description: "TW stock command",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "stock",
				Description: "股票名稱或代號，或輸入tw查詢加權指數",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionString,
			},
			{
				Name:        "args",
				Description: "額外參數",
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionString,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "三大法人買賣超日報",
						Value: "d",
					},
					{
						Name:  "三大法人買賣超月報",
						Value: "m",
					},
					{
						Name:  "股票資訊",
						Value: "i",
					},
				},
			},
		},
		Type: 1,
	})
	return
}
