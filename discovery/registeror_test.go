package discovery

import (
	"fmt"
	"strconv"
	"testing"
)

func TestRegister(t *testing.T) {
	pub := NewEtcdPublisher([]string{"http://127.0.0.1:2379"}, 10)
	id, _ := pub.lookupID("greeter")
	fmt.Println(id)

	prev := ""
	value := "1"
	if id > 0 {
		prev = strconv.Itoa(id)
		value = strconv.Itoa(id + 1)
	}
	b := pub.saveIDIndex("greeter", prev, value)
	fmt.Println(b)
}
