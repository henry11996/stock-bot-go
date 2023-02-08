package pkg

import (
	"fmt"
	"log"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/henry11996/fugle-golang/fugle"
	"github.com/shopspring/decimal"
)

func convertInfo(data fugle.Data) string {
	status := ""

	if data.Meta.IsSuspended {
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

	var currentPrice decimal.Decimal
	zero := decimal.NewFromInt(0)
	if data.Quote.Trade.Price.Equal(zero) {
		currentPrice = data.Quote.Trial.Price
	} else {
		currentPrice = data.Quote.Trade.Price
	}

	fivePricesText, totalUnitText := "", ""
	totalAskUnit, totalBidUnit := 0, 0
	for i := 0; i < 5; i++ {
		askPrice, askUnit, bidPrice, bidUnit := "", "", "", ""
		if len(data.Quote.Order.Asks) > i {
			bestasks := data.Quote.Order.Asks[i]
			askPrice = bestasks.Price.StringFixed(2)
			totalAskUnit += bestasks.Volume
			if askPrice == "0.00" {
				askPrice = "市價"
			}
			askUnit = strconv.Itoa(bestasks.Volume)
		}
		if len(data.Quote.Order.Bids) > i {
			bestbids := data.Quote.Order.Bids[i]
			bidPrice = bestbids.Price.StringFixed(2)
			totalBidUnit += bestbids.Volume
			if bidPrice == "0.00" {
				bidPrice = "市價"
			}
			bidUnit = strconv.Itoa(bestbids.Volume)
		}
		fivePricesText += fmt.Sprintf("%-5s %6s \\| %6s %5s\n", bidUnit, bidPrice, askPrice, askUnit)
	}
	totalUnitText += fmt.Sprintf("%-12v   %12v\n", totalBidUnit, totalAskUnit)

	if data.Meta.NameZhTw == "" {
		return "找不到此股票"
	} else if totalBidUnit == 0 && totalAskUnit == 0 {
		return fmt.Sprintf("``` %9s(%s)  %s \n"+
			"高 %4v\\ |低 %4v\\ |總 %5v\n"+
			"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
			"            %v         \n"+
			"        %v %2.2f%%      \n"+
			"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
			"```", data.Meta.NameZhTw, data.Info.SymbolID, status,
			data.Quote.PriceHigh.Price, data.Quote.PriceLow.Price, data.Quote.Total.TradeVolume,
			currentPrice.BigFloat(), data.Quote.Change.StringFixed(2), data.Quote.ChangePercent,
		)
	} else {
		return fmt.Sprintf("``` %9s(%s)  %s \n"+
			"高 %4v\\ |低 %4v\\ |總 %5v\n"+
			"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
			"            %v         \n"+
			" 買      %v %2.2f%%     賣\n"+
			"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
			"%s"+
			"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
			"%s"+
			"```", data.Meta.NameZhTw, data.Info.SymbolID, status,
			data.Quote.PriceHigh.Price, data.Quote.PriceLow.Price, data.Quote.Total.TradeVolume,
			currentPrice.BigFloat(), data.Quote.Change.StringFixed(2), data.Quote.ChangePercent,
			fivePricesText, totalUnitText,
		)
	}
}

func convert24TickerPrice(res *binance.PriceChangeStats) string {
	priceChangePercent, error := strconv.ParseFloat(res.PriceChangePercent, 64)
	if error != nil {
		log.Panic(error.Error())
	}
	return fmt.Sprintf("``` %13s 24hr \n"+
		"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
		"      %8s \n"+
		"   %8s %4.2f%% \n"+
		"\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\n"+
		"    高 %13s \n"+
		"    低 %13s \n"+
		"    總 %13s \n"+
		"```",
		res.Symbol,
		res.LastPrice,
		res.PriceChange,
		priceChangePercent,
		res.HighPrice,
		res.LowPrice,
		res.Volume,
	)
}
