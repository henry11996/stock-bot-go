package main

import (
	"log"
	"time"

	"github.com/telegram-go-stock-bot/plot"
)

func Route(command string, args []string, c chan interface{}) {
	var Now = time.Now().In(Loc)

	fugle := Initfugle()

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

	cmdType, param1 := "", ""
	if len(args) > 0 {
		cmdType = args[0]
	}
	if param1 != "" {
		param1 = args[1]
	}

	var err error
	switch command {
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
			lp, _ := getDayTotalLegalPerson(t)
			c <- lp.PrettyString()
		case "m":
			t := Now
			if param1 != "" {
				t, err = time.Parse("2006/01", param1)
				if err != nil {
					log.Panic("錯誤日期格式yyyy/mm")
				}
			}
			lp, _ := getMonthTotalLegalPerson(t)
			c <- lp.PrettyString()
		case "p":
			meta, _ := fugle.Meta("TWSE_SEM_INDEX_1", false)
			chart, _ := fugle.Chart("TWSE_SEM_INDEX_1", false)
			chart.Data.Meta = meta.Data.Meta
			chart.Data.Meta.NameZhTw = "加權指數"
			chart.Data.Info.SymbolID = "TWSE"
			png := plot.NewPlot(chart.Data, DefaultFont)
			c <- png
		default:
			meta, _ := fugle.Meta("TWSE_SEM_INDEX_1", false)
			quote, _ := fugle.Quote("TWSE_SEM_INDEX_1", false)
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
		case "p":
			meta, _ := fugle.Meta(stockId, false)
			chart, _ := fugle.Chart(stockId, false)
			chart.Data.Meta = meta.Data.Meta
			png := plot.NewPlot(chart.Data, DefaultFont)
			c <- png
		default:
			meta, _ := fugle.Meta(stockId, false)
			quote, _ := fugle.Quote(stockId, false)
			meta.Data.Quote = quote.Data.Quote
			c <- convertQuote(meta.Data)
		}
	}
}
