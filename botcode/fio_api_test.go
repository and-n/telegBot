package botcode

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_getBalance(t *testing.T) {
	var key = `{"accountStatement":"value"}`

	var balance S
	json.Unmarshal([]byte(key), &balance)

	fmt.Printf("statement\n %+v", balance.AccountStatement)

}

type S struct {
	AccountStatement string
}
