package botcode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const fioApi = "https://www.fio.cz/ib_api/rest"
const format = "transactions.json"

func getBalance(key string) string {
	tomorrow := time.Now().Add(time.Hour * 24).Format("2006-01-02")

	requestURL := fmt.Sprintf("%s/%s", fioApi, "periods/"+key+"/"+tomorrow+"/"+tomorrow+"/"+format)
	println(requestURL)
	res, err := http.Get(requestURL)
	if err != nil || res.StatusCode != 200 {
		fmt.Printf("error making http request: %s\n", err)
		return "error"
	}

	resBody, _ := ioutil.ReadAll(res.Body)

	var balance Balance
	json.Unmarshal([]byte(resBody), &balance)

	fmt.Printf("client: got response!\n %s", resBody)

	p := message.NewPrinter(language.English)

	return p.Sprintf("%.2f", balance.AccountStatement.Info.ClosingBalance) + " " + balance.AccountStatement.Info.Currency
}

type Balance struct {
	AccountStatement Statement
}

type Statement struct {
	Info            Info
	TransactionList string
}

type Info struct {
	Currency       string
	ClosingBalance float32
}
