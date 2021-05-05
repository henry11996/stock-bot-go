package main

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/RainrainWu/fugle-realtime-go/client"
)

func convertByTemplate(templa string, data client.FugleAPIData) (string, error) {
	t := template.New(templa + ".html")

	var err error
	t, err = t.ParseFiles("assets/html/" + templa + ".html")
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func convertQuote(data client.FugleAPIData) string {
	var status string
	if data.Quote.Istrial {
		status = "試搓中"
	} else if data.Quote.Iscurbingrise {
		status = "緩漲試搓"
	} else if data.Quote.Iscurbingfall {
		status = "緩跌試搓"
	} else if data.Quote.Isclosed {
		status = "已收盤"
	} else if data.Quote.Ishalting {
		status = "暫停交易"
	} else {
		status = "正常交易"
	}

	var bestPrices string
	for i, bestask := range data.Quote.Order.Bestasks {
		for j, _ := range data.Quote.Order.Bestbids {
			bid := data.Quote.Order.Bestbids[len(data.Quote.Order.Bestbids)-1-j]
			if i == j {
				bestPrices += fmt.Sprintf("%4d %5d \\| %4d %5d\n", bid.Price.IntPart(), bid.Unit.IntPart(), bestask.Price.IntPart(), bestask.Unit.IntPart())
			}
		}
	}

	return fmt.Sprintf("``` %4s \\- %s \\- %s\n"+
		"高 %v \\| 低 %v \\| 總 %v\n"+
		"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
		"    委買   %v   委賣\n"+
		"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
		"%s```", data.Meta.Namezhtw, data.Info.SymbolID, status,
		data.Quote.PriceHigh.Price, data.Quote.PriceLow.Price, data.Quote.Total.Unit,
		data.Quote.Trade.Price,
		bestPrices,
	)
}
