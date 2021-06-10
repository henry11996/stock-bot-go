package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type QueryResponse struct {
	ResultSet struct {
		Query  string
		Result []struct {
			Symbol   string
			Name     string
			Exch     string
			Type     string
			ExchDisp string
			TypeDisp string
		}
	}
}

type LegalPersonResponse struct {
	Stat   string
	Date   string
	Title  string
	Fields []string
	Data   [][]string
}

type LegalPerson struct {
	Title      string
	Date       string
	StockId    string
	StockName  string
	Foreign    LegalPersonTransaction
	Investment LegalPersonTransaction
	Dealer     LegalPersonTransaction
	Total      LegalPersonTransaction
}

type LegalPersonTransaction struct {
	Buy   int
	Sell  int
	Total int
}

func query(stockName string) (string, error) {
	template := "https://tw.stock.yahoo.com/stock_ms/_td-stock/api/resource/AutocompleteService;query=%s"
	r, err := http.Get(fmt.Sprintf(template, stockName))
	if err != nil {
		log.Print(err)
	}
	defer r.Body.Close()
	response := &QueryResponse{}
	json.NewDecoder(r.Body).Decode(response)

	stockId := response.getStockId()
	if stockId == "" {
		return stockName, nil
	} else {
		return stockId, nil
	}
}

func (res *QueryResponse) getStockId() string {
	if len(res.ResultSet.Result) > 0 {
		return strings.Split(res.ResultSet.Result[0].Symbol, ".")[0]
	}
	return ""
}

func getMonthLegalPersons(date string) LegalPersonResponse {
	url := "https://www.twse.com.tw/fund/TWT47U?response=json&selectType=ALL&date=%s"
	r, err := http.Get(fmt.Sprintf(url, date))
	if err != nil {
		log.Print(err)
	}
	defer r.Body.Close()
	response := &LegalPersonResponse{}
	json.NewDecoder(r.Body).Decode(response)
	return *response
}

func getDateLegalPersons(date string) LegalPersonResponse {
	url := "https://www.twse.com.tw/fund/T86?response=json&selectType=ALL&date=%s"
	r, err := http.Get(fmt.Sprintf(url, date))
	if err != nil {
		log.Print(err)
	}
	defer r.Body.Close()
	response := &LegalPersonResponse{}
	json.NewDecoder(r.Body).Decode(response)
	return *response
}

func (res *LegalPersonResponse) getByStock(stockId string, stockName string) LegalPerson {
	if len(res.Data) > 0 {
		for i := 0; i < len(res.Data); i++ {
			stockData := res.Data[i]
			if len(stockData) == len(res.Fields) && len(stockData) == 19 {
				if strings.Trim(stockData[0], " ") == stockId || strings.Trim(stockData[1], " ") == stockName {
					return res.NewLegalPerson(stockData)
				}
			}
		}
	}
	return LegalPerson{}
}

func (res *LegalPersonResponse) NewLegalPerson(stockData []string) LegalPerson {
	formatNumber := func(s string) int {
		i, err := strconv.Atoi(strings.ReplaceAll(strings.Trim(s, " "), ",", ""))
		if err != nil {
			panic(err)
		}
		return i
	}

	foreignBuy := formatNumber(stockData[2]) + formatNumber(stockData[5])
	foreignSell := formatNumber(stockData[3]) + formatNumber(stockData[6])
	foreignTotal := formatNumber(stockData[4]) + formatNumber(stockData[7])

	investmentBuy := formatNumber(stockData[8])
	investmentSell := formatNumber(stockData[9])
	investmentTotal := formatNumber(stockData[10])

	dealerBuy := formatNumber(stockData[12]) + formatNumber(stockData[15])
	dealerSell := formatNumber(stockData[13]) + formatNumber(stockData[16])
	dealerTotal := formatNumber(stockData[11])

	total := formatNumber(stockData[18])

	legalPerson := &LegalPerson{
		Date:      res.Date,
		StockId:   strings.Trim(stockData[0], " "),
		StockName: strings.Trim(stockData[1], " "),
		Title:     res.Title,
		Foreign: LegalPersonTransaction{
			Buy:   foreignBuy,
			Sell:  foreignSell,
			Total: foreignTotal,
		},
		Investment: LegalPersonTransaction{
			Buy:   investmentBuy,
			Sell:  investmentSell,
			Total: investmentTotal,
		},
		Dealer: LegalPersonTransaction{
			Buy:   dealerBuy,
			Sell:  dealerSell,
			Total: dealerTotal,
		},
		Total: LegalPersonTransaction{
			Buy:   foreignBuy + investmentBuy + dealerBuy,
			Sell:  foreignSell + investmentSell + dealerSell,
			Total: total,
		},
	}
	return *legalPerson
}
