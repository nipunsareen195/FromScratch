package temp

import (
	"fmt"
	"testing"
)

func TestStart(t *testing.T) {
	tx := Transaction{
		From:  "add1",
		To:    "add2",
		Value: 100,
	}

	by := EncodeToBytes(tx)

	fmt.Println(by)

	p := DecodeToPerson(by)

	fmt.Println(p)

	// reqBodyBytes := new(bytes.Buffer)
	// json.NewEncoder(reqBodyBytes).Encode(tx)

	// by := reqBodyBytes.Bytes() // this is the []byte

	// fmt.Println(by)

}
