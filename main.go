package main

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/freetype/truetype"
)

var Loc, _ = time.LoadLocation("Asia/Taipei")
var DefaultFont = readFont()

func main() {
	InitEnv()
	bot, updates := Boot()
	InitCache()
	InitSchedule()
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s(%v)] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

		c := make(chan interface{})

		args := strings.Split(update.Message.Text, " ")
		if len(args) > 1 {
			args = args[1:]
		} else {
			args = make([]string, 0)
		}

		go run(strings.Split(update.Message.Text, " ")[0][1:], args, c)

		message := <-c
		switch v := message.(type) {
		case string:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = v
			msg.ParseMode = "MarkdownV2"
			_, err := bot.Send(msg)
			if err != nil {
				log.Print(err)
			}
		case []byte:
			msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{
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

func run(command string, args []string, c chan interface{}) {
	var Now = time.Now().In(Loc)

	fugle := fugleInit()

	defer func() {
		if r := recover(); r != nil {
			err := ""
			switch x := r.(type) {
			case string:
				err = x
			case error:
				err = x.Error()
			}
			c <- "無法取得```" + command + "```，" + err
		}
	}()

	stockId, _ := query(command)
	var err error
	if command == "tw" {
		if len(args) > 0 && args[0] == "d" {
			t := Now
			if len(args) > 1 {
				t, err = time.Parse("2006/01/02", args[1])
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm/dd")
				}
			}
			lp, err := getDayTotalLegalPerson(t)
			if err != nil {
				log.Panic(err)
			}
			c <- lp.PrettyString()
		} else if len(args) > 0 && args[0] == "m" {
			t := Now
			if len(args) > 1 {
				t, err = time.Parse("2006/01", args[1])
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm")
				}
			}
			lp, err := getMonthTotalLegalPerson(t)
			if err != nil {
				log.Panic(err)
			}
			c <- lp.PrettyString()
		} else if len(args) > 0 && args[0] == "p" {
			meta, _ := fugle.Meta("TWSE_SEM_INDEX_1", false)
			chart, _ := fugle.Chart("TWSE_SEM_INDEX_1", false)
			chart.Data.Meta = meta.Data.Meta
			chart.Data.Meta.NameZhTw = "加權指數"
			chart.Data.Info.SymbolID = "TWSE"
			png := newPlot(chart.Data)
			c <- png
		} else {
			meta, _ := fugle.Meta("TWSE_SEM_INDEX_1", false)
			quote, _ := fugle.Quote("TWSE_SEM_INDEX_1", false)
			meta.Data.Meta.NameZhTw = "加權指數"
			meta.Data.Info.SymbolID = "TWSE"
			meta.Data.Quote = quote.Data.Quote
			c <- convertQuote(meta.Data)
		}
	} else {
		if len(args) > 0 && args[0] == "i" {
			meta, _ := fugle.Meta(stockId, false)
			c <- convertInfo(meta.Data)
		} else if len(args) > 0 && args[0] == "d" {
			t := Now
			if len(args) > 1 {
				t, err = time.Parse("2006/01/02", args[1])
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm/dd")
				}
			}
			lp, _ := getDayLegalPersons(t)

			c <- lp.FindStock(stockId, command).PrettyString(lp.Title)
		} else if len(args) > 0 && args[0] == "m" {
			t := Now
			if len(args) > 1 {
				t, err = time.Parse("2006/01", args[1])
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm")
				}
			}
			lp, _ := getMonthLegalPersons(t)

			c <- lp.FindStock(stockId, command).PrettyString(lp.Title)
		} else if len(args) > 0 && args[0] == "p" {
			meta, _ := fugle.Meta(stockId, false)
			res, _ := fugle.Chart(stockId, false)
			res.Data.Meta = meta.Data.Meta
			png := newPlot(res.Data)
			c <- png
		} else {
			meta, _ := fugle.Meta(stockId, false)
			quote, _ := fugle.Quote(stockId, false)
			meta.Data.Quote = quote.Data.Quote
			c <- convertQuote(meta.Data)
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
