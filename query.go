package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func query(stockName string, target interface{}) error {
	template := "https://tw.stock.yahoo.com/stock_ms/_td-stock/api/resource/AutocompleteService;query=%s"
	r, err := http.Get(fmt.Sprintf(template, stockName))
	if err != nil {
		// handle error
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
