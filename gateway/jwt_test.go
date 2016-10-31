package gateway

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGen(t *testing.T) {
	ret, err := Generate(11, "silentred", "email")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ret)

	userClaims, err := Verify(ret)
	if err != nil {
		fmt.Println(err)
	}
	b, _ := json.Marshal(userClaims)
	fmt.Println(string(b))
}
