package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/telegram-go-stock-bot/twse"
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

func getDayTotalLegalPerson(date time.Time) (*twse.LegalPersonTotal, error) {
	cacheKey := "day_total_legal_persons_"
	if x, found := Cache.Get(cacheKey + date.Format("20060102")); found {
		return x.(*twse.LegalPersonTotal), nil
	}
	totalPerson, err := twse.DayLegalPersonTotal(date)
	if err != nil {
		return &twse.LegalPersonTotal{}, err
	}
	Cache.Set(cacheKey+totalPerson.Date, totalPerson, cache.NoExpiration)
	return totalPerson, nil
}

func getMonthTotalLegalPerson(date time.Time) (*twse.LegalPersonTotal, error) {
	cacheKey := "month_total_legal_persons_"
	if x, found := Cache.Get(cacheKey + date.Format("200601") + "01"); found {
		return x.(*twse.LegalPersonTotal), nil
	}
	totalPerson, err := twse.MonthLegalPersonTotal(date)
	if err != nil {
		return &twse.LegalPersonTotal{}, err
	}
	Cache.Set(cacheKey+totalPerson.Date, totalPerson, cache.NoExpiration)
	return totalPerson, nil
}

func getMonthLegalPersons(date time.Time) (*twse.LegalPersonStocks, error) {
	cacheKey := "month_legal_persons_"
	if x, found := Cache.Get(cacheKey + date.Format("200601") + "01"); found {
		return x.(*twse.LegalPersonStocks), nil
	}
	legalPersons, err := twse.MonthLegalPersons(date)
	if err != nil {
		return &twse.LegalPersonStocks{}, err
	}
	log.Print(legalPersons.Date)
	Cache.Set(cacheKey+legalPersons.Date, legalPersons, cache.NoExpiration)
	return legalPersons, nil
}

func getDayLegalPersons(date time.Time) (*twse.LegalPersonStocks, error) {
	cacheKey := "day_legal_persons_"
	if x, found := Cache.Get(cacheKey + date.Format("20060102")); found {
		return x.(*twse.LegalPersonStocks), nil
	}
	legalPersons, err := twse.DayLegalPersons(date)
	if err != nil {
		return &twse.LegalPersonStocks{}, err
	}
	otclegalPersons, err := twse.DayOTCLegalPersons(date)
	if err != nil {
		return &twse.LegalPersonStocks{}, err
	}
	legalPersons.Stocks = append(legalPersons.Stocks, otclegalPersons.Stocks...)
	Cache.Set(cacheKey+legalPersons.Date, legalPersons, cache.NoExpiration)
	return legalPersons, nil
}
