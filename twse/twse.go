package twse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LegalPersonResponse struct {
	Stat   string
	Date   string
	Title  string
	Fields []string
	Data   [][]string
}

type LegalPersonTotal struct {
	Stat       string
	Date       string
	Title      string
	Foreign    LegalPersonTransaction
	Investment LegalPersonTransaction
	Dealer     LegalPersonTransaction
	Total      LegalPersonTransaction
}

type LegalPersonStocks struct {
	Stat   string
	Date   string
	Title  string
	Stocks []LegalPersonStock
}

type LegalPersonStock struct {
	Id         string
	Name       string
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

var Loc, _ = time.LoadLocation("Asia/Taipei")

func DayLegalPersonTotal(date time.Time) (*LegalPersonTotal, error) {
	var url string
	if date.Equal(time.Now().In(Loc)) {
		url = "https://www.twse.com.tw/fund/BFI82U?response=json&weekDate=&monthDate=&type=day"
	} else {
		url = fmt.Sprintf("https://www.twse.com.tw/fund/BFI82U?response=json&weekDate=&monthDate=&type=day&dayDate=%s", date.Format("20060102"))
	}
	response := &LegalPersonResponse{}
	legalPerson := &LegalPersonTotal{}
	err := request(url, response)
	if err != nil {
		return legalPerson, err
	}
	legalPerson = NewLegalPersonTotal(response.Date, response.Title, response.Data)
	return legalPerson, nil
}

func MonthLegalPersons(date time.Time) (*LegalPersonStocks, error) {
	var url string
	if date.Equal(time.Now().In(Loc)) {
		url = "https://www.twse.com.tw/fund/TWT47U?response=json&selectType=ALL"
	} else {
		url = fmt.Sprintf("https://www.twse.com.tw/fund/TWT47U?response=json&selectType=ALL&date=%s", date.Format("20060102"))
	}
	stocks := &LegalPersonStocks{}
	response := &LegalPersonResponse{}
	err := request(url, response)
	if err != nil {
		return stocks, err
	}
	stocks = &LegalPersonStocks{
		Title: response.Title,
		Stat:  response.Stat,
		Date:  response.Date,
	}
	for _, data := range response.Data {
		stocks.Stocks = append(stocks.Stocks, NewLegalPersonStock(data))
	}
	return stocks, nil
}

func DayLegalPersons(date time.Time) (*LegalPersonStocks, error) {
	var url string
	if date.Equal(time.Now().In(Loc)) {
		url = "https://www.twse.com.tw/fund/T86?response=json&selectType=ALL"
	} else {
		url = fmt.Sprintf("https://www.twse.com.tw/fund/T86?response=json&selectType=ALL&date=%s", date.Format("20060102"))
	}
	stocks := &LegalPersonStocks{}
	response := &LegalPersonResponse{}
	err := request(url, response)
	if err != nil {
		return stocks, err
	}
	stocks = &LegalPersonStocks{
		Title: response.Title,
		Stat:  response.Stat,
		Date:  response.Date,
	}
	for _, data := range response.Data {
		stocks.Stocks = append(stocks.Stocks, NewLegalPersonStock(data))
	}
	return stocks, nil
}

func request(url string, x interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(x)
}

func (lp *LegalPersonStocks) FindStock(stockId string, stockName string) *LegalPersonStock {
	for _, stock := range lp.Stocks {
		if stock.Id == stockId || stock.Name == stockName {
			return &stock
		}
	}
	return &LegalPersonStock{}
}

func NewLegalPersonStock(stockData []string) LegalPersonStock {
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

	legalPerson := &LegalPersonStock{
		Id:   strings.Trim(stockData[0], " "),
		Name: strings.Trim(stockData[1], " "),
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

func (lp *LegalPersonStock) PrettyString(title string) string {
	typeTitle := strings.Split(title, " ")[1]
	dateTitle := strings.Split(title, " ")[0]
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
		lp.Name, lp.Id,
		center(typeTitle, w),
		center(dateTitle, w),
		lp.Foreign.Buy/1000, lp.Foreign.Sell/1000, lp.Foreign.Total/1000,
		lp.Investment.Buy/1000, lp.Investment.Sell/1000, lp.Investment.Total/1000,
		lp.Dealer.Buy/1000, lp.Dealer.Sell/1000, lp.Dealer.Total/1000,
		lp.Total.Buy/1000, lp.Total.Sell/1000, lp.Total.Total/1000,
	)
}

func NewLegalPersonTotal(date string, title string, stockData [][]string) *LegalPersonTotal {
	formatNumber := func(s string) int {
		i, err := strconv.Atoi(strings.ReplaceAll(strings.Trim(s, " "), ",", ""))
		if err != nil {
			panic(err)
		}
		return i
	}

	foreignBuy := formatNumber(stockData[3][1]) + formatNumber(stockData[4][1])
	foreignSell := formatNumber(stockData[3][2]) + formatNumber(stockData[4][2])
	foreignTotal := formatNumber(stockData[3][3]) + formatNumber(stockData[4][3])

	investmentBuy := formatNumber(stockData[2][1])
	investmentSell := formatNumber(stockData[2][2])
	investmentTotal := formatNumber(stockData[2][3])

	dealerBuy := formatNumber(stockData[0][1]) + formatNumber(stockData[1][1])
	dealerSell := formatNumber(stockData[0][2]) + formatNumber(stockData[1][2])
	dealerTotal := formatNumber(stockData[0][3]) + formatNumber(stockData[1][3])

	total := formatNumber(stockData[5][3])

	legalPerson := &LegalPersonTotal{
		Date:  date,
		Title: title,
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
	return legalPerson
}

func (lp *LegalPersonTotal) PrettyString() string {
	typeTitle := strings.Split(lp.Title, " ")[1]
	dateTitle := strings.Split(lp.Title, " ")[0]
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
		float64(lp.Foreign.Buy)/100000000, float64(lp.Foreign.Sell)/100000000, float64(lp.Foreign.Total)/100000000,
		float64(lp.Investment.Buy)/100000000, float64(lp.Investment.Sell)/100000000, float64(lp.Investment.Total)/100000000,
		float64(lp.Dealer.Buy)/100000000, float64(lp.Dealer.Sell)/100000000, float64(lp.Dealer.Total)/100000000,
		float64(lp.Total.Buy)/100000000, float64(lp.Total.Sell)/100000000, float64(lp.Total.Total)/100000000,
	)
}
