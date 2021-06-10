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

func getLegalPersons(date string) LegalPersonResponse {
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
	var err error

	var foreignChianBuy, foreignBuy, foreignChianSell, foreignSell, foreignChianTotal, foreignTotal int
	foreignChianBuy, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[2], " "), ",", ""))
	foreignBuy, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[5], " "), ",", ""))
	foreignBuy += foreignChianBuy
	foreignChianSell, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[3], " "), ",", ""))
	foreignSell, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[6], " "), ",", ""))
	foreignSell += foreignChianSell
	foreignChianTotal, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[4], " "), ",", ""))
	foreignTotal, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[7], " "), ",", ""))
	foreignTotal += foreignChianTotal

	var investmentBuy, investmentSell, investmentTotal int
	investmentBuy, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[8], " "), ",", ""))
	investmentSell, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[9], " "), ",", ""))
	investmentTotal, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[10], " "), ",", ""))

	var dealerSelfBuy, dealerSelfSell, dealerBuy, dealerSell, dealerTotal int
	dealerSelfBuy, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[12], " "), ",", ""))
	dealerBuy, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[15], " "), ",", ""))
	dealerBuy += dealerSelfBuy
	dealerSelfSell, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[13], " "), ",", ""))
	dealerSell, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[16], " "), ",", ""))
	dealerSell += dealerSelfSell
	dealerTotal, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[11], " "), ",", ""))

	var total int
	total, err = strconv.Atoi(strings.ReplaceAll(strings.Trim(stockData[18], " "), ",", ""))

	if err != nil {
		panic(err)
	}

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
