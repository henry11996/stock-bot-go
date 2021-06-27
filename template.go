package main

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/henry11996/fugle-golang/fugle"
	"github.com/shopspring/decimal"
)

func convertInfo(data fugle.Data) string {
	status := ""

	if data.Meta.IsSuspended {
		status += "暫停買賣 "
	}
	if data.Meta.CanShortMargin && data.Meta.CanShortLend {
		status += "暫停買賣 "
	}
	if data.Meta.CanShortMargin && data.Meta.CanShortLend {
		status += "可融資券 "
	} else if data.Meta.CanShortMargin {
		status += "禁融券 "
	} else if data.Meta.CanShortLend {
		status += "禁融資 "
	} else {
		status += "禁融資券 "
	}

	if data.Meta.CanDayBuySell && data.Meta.CanDaySellBuy {
		status += "買賣現沖 "
	} else if data.Meta.CanDayBuySell {
		status += "現沖買 "
	} else if data.Meta.CanDaySellBuy {
		status += "現沖賣 "
	} else {
		status += "禁現沖 "
	}

	return fmt.Sprintf("```\n"+
		"%s(%s)\n"+
		"產業：%s\n"+
		"狀態：%s\n"+
		"現價：%s```\n",
		data.Meta.NameZhTw, data.Info.SymbolID,
		data.Meta.IndustryZhTw,
		status,
		data.Meta.PriceReference,
	)
}

func convertQuote(data fugle.Data) string {
	var status string
	if data.Quote.IsTrial {
		status = "試搓中"
	} else if data.Quote.IsCurbingRise {
		status = "緩漲試搓"
	} else if data.Quote.IsCurbingFall {
		status = "緩跌試搓"
	} else if data.Quote.IsClosed {
		//已收盤
		status = ""
	} else if data.Quote.IsHalting {
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
	percent = currentPirce.Sub(data.Meta.PriceReference).Div(data.Meta.PriceReference).Mul(hunded).BigFloat()
	minus = currentPirce.Sub(data.Meta.PriceReference).BigFloat()
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
				bidUnit = strconv.Itoa(bestbids.Unit)
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
					askUnit = strconv.Itoa(bestasks.Unit)
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
		"%s```", data.Meta.NameZhTw, data.Info.SymbolID, status,
		data.Quote.PriceHigh.Price, data.Quote.PriceLow.Price, data.Quote.Total.Unit,
		currentPirce.BigFloat(), minus, percent,
		bestPrices,
	)
}

func convertLegalPerson(legalPerson LegalPerson) string {
	if legalPerson.Title == "" {
		panic("找不到買賣超資料")
	}
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
