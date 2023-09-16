package botcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	ttlcache "github.com/jellydator/ttlcache/v3"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const fioApi = "https://www.fio.cz/ib_api/rest"
const format = "transactions.json"
const dateFormat = "2006-01-02"

var cacheBalance *ttlcache.Cache[string, AccountStatement]

func init() {
	cacheBalance = ttlcache.New(
		ttlcache.WithTTL[string, AccountStatement](30 * time.Minute),
	)

	go cacheBalance.Start()
}

func getBalance(key string) string {
	var data struct {
		AccountStatement AccountStatement `json:"accountStatement"`
	}
	var balance AccountStatement

	if cacheBalance != nil && cacheBalance.Get(key) != nil {
		fmt.Println("get from cache")
		balance = cacheBalance.Get(key).Value()
	} else {

		tomorrow := time.Now().Add(time.Hour * 24).Format(dateFormat)

		requestURL := fmt.Sprintf("%s/%s", fioApi, "periods/"+key+"/"+tomorrow+"/"+tomorrow+"/"+format)
		// println(requestURL)
		res, err := http.Get(requestURL)
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
		// fmt.Printf("client: got response!\n %s\n", resBody)
		// fmt.Println(balance)

		cacheBalance.Set(key, balance, time.Minute*5)
	}

	p := message.NewPrinter(language.English)

	return p.Sprintf("%.2f", balance.Info.ClosingBalance) + " " + balance.Info.Currency
}

func getSumByMonth(key string, month int) (Balance, error) {

	first, last, err := getfirstAndLastDayOfMonth(month)
	if err != nil {
		return Balance{}, err
	}

	requestURL := fmt.Sprintf("%s/%s", fioApi, "periods/"+key+"/"+first.Format(dateFormat)+"/"+last.Format(dateFormat)+"/"+format)

}

func getfirstAndLastDayOfMonth(month int) (first time.Time, last time.Time, err error) {
	if month <= 0 || month > 12 {
		return time.Time{}, time.Time{}, errors.New("Wrong month number. (1-12)")
	}

	today := time.Now()

	var firstDayOfMonth, lastDayOfMonth time.Time
	if month <= int(today.Month()) {

		firstDayOfMonth = time.Date(today.Year(), time.Month(month), 1, 0, 0, 1, 0, today.Location())
		lastDayOfMonth = time.Date(today.Year(), time.Month(month)+1, month+1, 0, 0, 1, 0, today.Location()).Add(time.Minute * -1)

	} else {
		firstDayOfMonth = time.Date(today.Year()-1, time.Month(month), 1, 0, 0, 1, 0, today.Location())
		// lastDayOfMonth =

	}
	return firstDayOfMonth, lastDayOfMonth, nil
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
	Value float64 `json:"column_1"`
	// Add your transaction fields here
}

type Balance struct {
	Income   float64
	Expenses float64
	Total    float64
}
