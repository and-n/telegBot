package botcode

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const fioApi = "https://www.fio.cz/ib_api/rest"
const format = "transactions.json"

func getBalance(key string) string {
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

		json.Unmarshal([]byte(resBody), &balance)
		// fmt.Printf("client: got response!\n %s", resBody)

		cacheBalance.Set(key, balance, time.Minute*5)
	}

	p := message.NewPrinter(language.English)

	return p.Sprintf("%.2f", balance.Info.ClosingBalance) + " " + balance.Info.Currency
}

type AccountStatement struct {
	Info            Info
	TransactionList string
}

type Info struct {
	Currency       string
	ClosingBalance float32
	AccountId      int64
	BankId         int16
	Iban           string
	Bic            string
}
