package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/RainrainWu/fugle-realtime-go/client"
	"github.com/shopspring/decimal"
)

func convertInfo(data client.FugleAPIData) string {
	status := ""

	if data.Meta.Issuspended {
		status += "暫停買賣 "
	}
	if data.Meta.Canshortmargin && data.Meta.Canshortlend {
		status += "暫停買賣 "
	}
	if data.Meta.Canshortmargin && data.Meta.Canshortlend {
		status += "可融資券 "
	} else if data.Meta.Canshortmargin {
		status += "禁融券 "
	} else if data.Meta.Canshortlend {
		status += "禁融資 "
	} else {
		status += "禁融資券 "
	}

	if data.Meta.Candaybuysell && data.Meta.Candaysellbuy {
		status += "買賣現沖 "
	} else if data.Meta.Candaybuysell {
		status += "現沖買 "
	} else if data.Meta.Candaysellbuy {
		status += "現沖賣 "
	} else {
		status += "禁現沖 "
	}

	return fmt.Sprintf("[%s\\(%s\\)](https://tw.stock.yahoo.com/q/bc?s=%s)\n"+
		"產業：%s\n"+
		"狀態：%s\n"+
		"現價：%s\n",
		data.Meta.Namezhtw, data.Info.SymbolID, data.Info.SymbolID,
		data.Meta.Industryzhtw,
		status,
		data.Meta.Pricereference,
	)
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
		//已收盤
		status = ""
	} else if data.Quote.Ishalting {
		status = "暫停交易"
	} else {
		//正常交易
		status = ""
	}

	var currentPirce decimal.Decimal
	zero := decimal.NewFromInt(0)
	if data.Quote.Trade.Price.Equal(zero) {
		currentPirce = data.Quote.Trial.Price
	} else {
		currentPirce = data.Quote.Trade.Price
	}

	var percent, minus *big.Float
	hunded := decimal.NewFromInt(100)
	percent = currentPirce.Sub(data.Meta.Pricereference).Div(data.Meta.Pricereference).Mul(hunded).BigFloat()
	minus = currentPirce.Sub(data.Meta.Pricereference).BigFloat()
	var bestPrices string

	if len(data.Quote.Order.Bestbids) > 0 || len(data.Quote.Order.Bestasks) > 0 {
		for i := 0; i < 5; i++ {
			bidPrice := ""
			bidUnit := ""
			if len(data.Quote.Order.Bestbids) > i {
				bestbids := data.Quote.Order.Bestbids[len(data.Quote.Order.Bestbids)-1-i]
				bidPrice = bestbids.Price.StringFixed(2)
				if bidPrice == "0.00" {
					bidPrice = "市價"
				}
				bidUnit = bestbids.Unit.String()
			}
			for j := 0; j < 5; j++ {
				askPrice := ""
				askUnit := ""
				if len(data.Quote.Order.Bestasks) > j {
					bestasks := data.Quote.Order.Bestasks[j]
					askPrice = bestasks.Price.StringFixed(2)
					if askPrice == "0.00" {
						askPrice = "市價"
					}
					askUnit = bestasks.Unit.String()
				}
				if i == j {
					bestPrices += fmt.Sprintf("%6s %5s \\| %6s %5s\n", bidPrice, bidUnit, askPrice, askUnit)
				}
			}
		}
	} else {
		bestPrices = ""
	}

	return fmt.Sprintf("``` %9s(%s)  %s \n"+
		"高 %4v \\| 低 %4v \\| 總 %5v\n"+
		"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
		"            %v         \n"+
		"    買   %2.2f %2.2f%%   賣\n"+
		"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
		"%s```", data.Meta.Namezhtw, data.Info.SymbolID, status,
		data.Quote.PriceHigh.Price, data.Quote.PriceLow.Price, data.Quote.Total.Unit,
		currentPirce.BigFloat(), minus, percent,
		bestPrices,
	)
}

func convertLegalPerson(legalPerson LegalPerson) string {
	typeTitle := strings.Split(legalPerson.Title, " ")[1]
	dateTitle := strings.Split(legalPerson.Title, " ")[0]
	w := 22

	center := func(s string, w int) string {
		return fmt.Sprintf("%*s", w/2, s[:len(s)/2]) + fmt.Sprintf("%*s", -w/2, s[len(s)/2:])
	}

	return fmt.Sprintf("```\n"+
		"  %9s(%s)\n"+
		"---------------------------\n"+
		"%s\n"+
		" %s\n"+
		"---------------------------\n"+
		"(張)   買  |   賣  |    總\n"+
		"外資 %6v| %6v| %6v\n"+
		"投信 %6v| %6v| %6v\n"+
		"自營 %6v| %6v| %6v\n"+
		"總共 %6v| %6v| %6v\n"+
		"```",
		legalPerson.StockName, legalPerson.StockId,
		center(typeTitle, w),
		center(dateTitle, w),
		legalPerson.Foreign.Buy/1000, legalPerson.Foreign.Sell/1000, legalPerson.Foreign.Total/1000,
		legalPerson.Investment.Buy/1000, legalPerson.Investment.Sell/1000, legalPerson.Investment.Total/1000,
		legalPerson.Dealer.Buy/1000, legalPerson.Dealer.Sell/1000, legalPerson.Dealer.Total/1000,
		legalPerson.Total.Buy/1000, legalPerson.Total.Sell/1000, legalPerson.Total.Total/1000,
	)
}

func convertTotalLegalPerson(legalPerson LegalPerson) string {
	typeTitle := strings.Split(legalPerson.Title, " ")[1]
	dateTitle := strings.Split(legalPerson.Title, " ")[0]
	w := 24

	center := func(s string, w int) string {
		return fmt.Sprintf("%*s", w/2, s[:len(s)/2]) + fmt.Sprintf("%*s", -w/2, s[len(s)/2:])
	}

	return fmt.Sprintf("```\n"+
		"%s\n"+
		"   %s\n"+
		"----------------------------\n"+
		"(億)   買   |  賣   |   總\n"+
		"外資 %7.2f|%7.2f|%7.2f\n"+
		"投信 %7.2f|%7.2f|%7.2f\n"+
		"自營 %7.2f|%7.2f|%7.2f\n"+
		"總共 %7.2f|%7.2f|%7.2f\n"+
		"```",
		center(typeTitle, w),
		center(dateTitle, w),
		float64(legalPerson.Foreign.Buy)/100000000, float64(legalPerson.Foreign.Sell)/100000000, float64(legalPerson.Foreign.Total)/100000000,
		float64(legalPerson.Investment.Buy)/100000000, float64(legalPerson.Investment.Sell)/100000000, float64(legalPerson.Investment.Total)/100000000,
		float64(legalPerson.Dealer.Buy)/100000000, float64(legalPerson.Dealer.Sell)/100000000, float64(legalPerson.Dealer.Total)/100000000,
		float64(legalPerson.Total.Buy)/100000000, float64(legalPerson.Total.Sell)/100000000, float64(legalPerson.Total.Total)/100000000,
	)
}
