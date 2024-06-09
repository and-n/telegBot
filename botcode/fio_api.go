package botcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	ttlcache "github.com/jellydator/ttlcache/v3"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const fioApi = "https://fioapi.fio.cz/v1/rest"
const format = "transactions.json"
const dateFormat = "2006-01-02"

var cacheStatement *ttlcache.Cache[string, AccountStatement]
var cacheMonthBalance *ttlcache.Cache[string, Balance]

func init() {
	cacheStatement = ttlcache.New(
		ttlcache.WithTTL[string, AccountStatement](30 * time.Minute),
	)
	cacheMonthBalance = ttlcache.New(
		ttlcache.WithTTL[string, Balance](24 * time.Hour),
	)

	go cacheStatement.Start()
	go cacheMonthBalance.Start()
}

func getBalance(key string) string {
	var data struct {
		AccountStatement AccountStatement `json:"accountStatement"`
	}
	var balance AccountStatement

	if cacheStatement != nil && cacheStatement.Get(key) != nil {
		fmt.Println("get from cache")
		balance = cacheStatement.Get(key).Value()
	} else {

		tomorrow := time.Now().Add(time.Hour * 24).Format(dateFormat)

		requestURL := fmt.Sprintf("%s/%s", fioApi, "periods/"+key+"/"+tomorrow+"/"+tomorrow+"/"+format)
		// println(requestURL)

		res, err := getRequest(requestURL)
		if err != nil || res.StatusCode != 200 {
			fmt.Printf("error making http request: %s\n", err)
			return "error"
		}

		resBody, _ := io.ReadAll(res.Body)

		e := json.Unmarshal([]byte(resBody), &data)
		if e != nil {
			fmt.Println("Error:", e)
			return ""
		}
		balance = data.AccountStatement
		// fmt.Println(balance)

		cacheStatement.Set(key, balance, time.Minute*5)
	}

	p := message.NewPrinter(language.English)

	return p.Sprintf("%.2f %s", balance.Info.ClosingBalance, balance.Info.Currency)
}

func getSumByMonthAsString(key string, month int) (string, error) {
	val, err := getSumByMonth(key, month)
	if err != nil {
		return "", err
	}
	p := message.NewPrinter(language.English)
	return p.Sprintf("Income: %.2f\nExpenses: %.2f\nTotal: %.2f", val.Income, val.Expenses, val.Total), err
}

func getSumByMonth(key string, month int) (Balance, error) {
	var mBalance Balance
	if cacheStatement != nil && cacheMonthBalance.Get(key+"_M"+strconv.Itoa(month)) != nil {
		fmt.Println("get from cache")
		mBalance = cacheMonthBalance.Get(key + "_M" + strconv.Itoa(month)).Value()
	} else {

		first, last, err := getFirstAndLastDayOfMonth(month)
		if err != nil {
			return Balance{}, err
		}

		requestURL := fmt.Sprintf("%s/%s", fioApi, "periods/"+key+"/"+first.Format(dateFormat)+"/"+last.Format(dateFormat)+"/"+format)

		res, err := getRequest(requestURL)
		if err != nil || res.StatusCode != 200 {
			if res.StatusCode == 422 {
				fmt.Printf("error making http request: %d %s\n", res.StatusCode, err)
				return Balance{}, errors.New("you can't get data for more than 90 days")
			} else {
				fmt.Printf("error making http request: %d %s\n", res.StatusCode, err)
				return Balance{}, errors.New("error making http request")
			}
		}

		var data struct {
			AccountStatement AccountStatement `json:"accountStatement"`
		}
		resBody, _ := io.ReadAll(res.Body)

		e := json.Unmarshal([]byte(resBody), &data)
		if e != nil {
			fmt.Println("Error:", e)
			if e, ok := err.(*json.SyntaxError); ok {
				fmt.Printf("syntax error at byte offset %d", e.Offset)
			}
			return Balance{}, e
		}
		mBalance = Balance{}
		for _, val := range data.AccountStatement.TransactionList.Transaction {
			if val.Column1.Value > 0 {
				mBalance.Income += val.Column1.Value
			} else {
				mBalance.Expenses += val.Column1.Value
			}
			mBalance.Total += val.Column1.Value
		}
		cacheMonthBalance.Set(key+"_M"+strconv.Itoa(month), mBalance, ttlcache.DefaultTTL)
	}
	return mBalance, nil
}

func getFirstAndLastDayOfMonth(month int) (first time.Time, last time.Time, err error) {
	if month <= 0 || month > 12 {
		return time.Time{}, time.Time{}, errors.New("wrong month number " + strconv.Itoa(month))
	}

	today := time.Now()

	var firstDayOfMonth, lastDayOfMonth time.Time
	if month <= int(today.Month()) {
		firstDayOfMonth = time.Date(today.Year(), time.Month(month), 1, 0, 0, 1, 0, today.Location())

		if month == 12 {
			lastDayOfMonth = time.Date(today.Year(), time.Month(month), 31, 0, 0, 1, 0, today.Location())
		} else {
			lastDayOfMonth = time.Date(today.Year(), time.Month(month)+1, 1, 0, 0, 1, 0, today.Location()).Add(time.Minute * -1)
		}

	} else {
		firstDayOfMonth = time.Date(today.Year()-1, time.Month(month), 1, 0, 0, 1, 0, today.Location())

		if month == 12 {
			lastDayOfMonth = time.Date(today.Year()-1, time.Month(month), 31, 0, 0, 1, 0, today.Location())
		} else {
			lastDayOfMonth = time.Date(today.Year()-1, time.Month(month)+1, 1, 0, 0, 1, 0, today.Location()).Add(time.Minute * -1)
		}

	}
	return firstDayOfMonth, lastDayOfMonth, nil
}

func getRequest(requestURL string) (*http.Response, error) {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.ExpectContinueTimeout = 10 * time.Second
	tr.DisableKeepAlives = true
	tr.IdleConnTimeout = 10 * time.Second

	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Please stop banning me")

	res, err := client.Do(req)
	return res, err
}

type AccountStatement struct {
	Info            Info
	TransactionList TransactionList
}

type Info struct {
	Currency       string
	ClosingBalance float64
	AccountId      string
	BankId         string
	Iban           string
	Bic            string
}

type TransactionList struct {
	Transaction []Transaction `json:"transaction"`
}

type Transaction struct {
	Column1 Column1 `json:"column1"`
	// Add your transaction fields here
}

type Column1 struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Id    int     `json:"id"`
}

type Balance struct {
	Income   float64
	Expenses float64
	Total    float64
}
