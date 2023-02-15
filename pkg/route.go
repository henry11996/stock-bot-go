package pkg

import (
	"context"
	"log"
	"strings"
	"time"

	myBinance "github.com/stock-bot-go/pkg/service/binance"
)

var (
	Loc, _ = time.LoadLocation("Asia/Taipei")
)

type RouteFunc func(command string, args []string, c chan string)

func Route(command string, args []string, c chan string) {
	var Now = time.Now().In(Loc)

	fugle := InitFugle()

	binance := myBinance.GetInstance()

	defer func() {
		if r := recover(); r != nil {
			err := ""
			switch x := r.(type) {
			case string:
				err = x
			case error:
				err = x.Error()
			}
			c <- "執行```" + command + "```錯誤，" + err
		}
	}()

	stockId, _ := query(command)

	cmdType, param1 := "", ""
	if len(args) > 0 {
		cmdType = args[0]
	}
	if len(args) > 1 {
		param1 = args[1]
	}

	var err error
	switch command {
	case "e":
		cmdType = strings.ToUpper(cmdType)
		if param1 == "" {
			param1 = "USDT"
		}
		param1 = strings.ToUpper(param1)
		res, err := binance.NewListPriceChangeStatsService().Symbol(cmdType + param1).Do(context.Background())
		if err != nil {
			log.Println(err)
			panic("無法取得```" + cmdType + "```")
		}
		res[0].Symbol = map[string]string{
			"BTC": "₿",
			"ETH": "⟠",
		}[cmdType] + " " + res[0].Symbol
		c <- convert24TickerPrice(res[0])
	case "tw":
		switch cmdType {
		case "d":
			t := Now
			if param1 != "" {
				t, err = time.Parse("2006/01/02", param1)
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm/dd")
				}
			}
			lp, err := getDayTotalLegalPerson(t)
			if err != nil {
				log.Println(err)
				panic("無法取得```" + command + "```")
			}
			c <- lp.PrettyString()
		case "m":
			t := Now
			if param1 != "" {
				t, err = time.Parse("2006/01", param1)
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm")
				}
			}
			lp, err := getMonthTotalLegalPerson(t)
			if err != nil {
				log.Println(err)
				panic("無法取得```" + command + "```")
			}
			c <- lp.PrettyString()
		default:
			meta, _ := fugle.Meta("IX0001", false)
			quote, _ := fugle.Quote("IX0001", false)
			meta.Data.Meta.NameZhTw = "加權指數"
			meta.Data.Info.SymbolID = "TWSE"
			meta.Data.Quote = quote.Data.Quote
			c <- convertQuote(meta.Data)
		}
	default:
		switch cmdType {
		case "i":
			meta, _ := fugle.Meta(stockId, false)
			c <- convertInfo(meta.Data)
		case "d":
			t := Now
			if param1 != "" {
				t, err = time.Parse("2006/01/02", param1)
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm/dd")
				}
			}
			lp, _ := getDayLegalPersons(t)
			c <- lp.FindStock(stockId, command).PrettyString(lp.Title)
		case "m":
			t := Now
			if param1 != "" {
				t, err = time.Parse("2006/01", param1)
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm")
				}
			}
			lp, _ := getMonthLegalPersons(t)
			c <- lp.FindStock(stockId, command).PrettyString(lp.Title)
		default:
			meta, _ := fugle.Meta(stockId, false)
			quote, _ := fugle.Quote(stockId, false)
			meta.Data.Quote = quote.Data.Quote
			c <- convertQuote(meta.Data)
		}
	}
}
