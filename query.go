package main

import (
	"encoding/json"
	"fmt"
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

func query(stockName string) (string, error) {
	template := "https://tw.stock.yahoo.com/stock_ms/_td-stock/api/resource/AutocompleteService;query=%s"
	r, err := http.Get(fmt.Sprintf(template, stockName))
	if err != nil {
		// handle error
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

// func legalPerson(stcokId string) string {
// 	url := "https://api.tej.com.tw/api/datatables/TRAIL/TATINST1.json?api_key=%s"

// }
