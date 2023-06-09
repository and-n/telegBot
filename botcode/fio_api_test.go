package botcode

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_getBalance(t *testing.T) {
	var key = `{"accountStatement":{"info":{"currency":"CZK"}}}`

	var balance AccountStatement
	json.Unmarshal([]byte(key), &balance)

	fmt.Printf("statement\n %+v\n", balance.Info)

}


type S struct {
	AccountStatement string
}
