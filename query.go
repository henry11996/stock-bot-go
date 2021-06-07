package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	Date      string
	StockId   string
	StockName string
	Buy       string
	Sell      string
	Total     string
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
	url := "https://www.twse.com.tw/fund/TWT38U?response=json&date=%s"
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
			if len(stockData) == len(res.Fields) && len(stockData) == 12 {
				if strings.Trim(stockData[1], " ") == stockId || strings.Trim(stockData[2], " ") == stockName {
					legalPerson := &LegalPerson{
						Date:      res.Date,
						StockId:   strings.Trim(stockData[1], " "),
						StockName: strings.Trim(stockData[2], " "),
						Buy:       strings.Trim(stockData[3], " "),
						Sell:      strings.Trim(stockData[4], " "),
						Total:     strings.Trim(stockData[5], " "),
					}
					return *legalPerson
				}
			}
		}
	}
	return LegalPerson{}
}
