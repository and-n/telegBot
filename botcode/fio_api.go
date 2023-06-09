package botcode

import (
	"encoding/json"
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
	if cacheBalance.Get(key) != nil {
		fmt.Println("get from cache")
		balance = cacheBalance.Get(key).Value()
	} else {

		tomorrow := time.Now().Add(time.Hour * 24).Format("2006-01-02")

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
	// Add your transaction fields here
}