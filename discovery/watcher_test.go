package discovery

import (
	"fmt"
	"testing"
)

func TestDiff(t *testing.T) {
	a := []string{"1", "2", "4"}
	b := []string{"5", "2"}
	c := diffString(a, b)
	fmt.Println(c)
}

func TestIntersect(t *testing.T) {
	a := []string{"1", "2", "4"}
	b := []string{"5", "2", "2"}
	c := intersectString(a, b)
	fmt.Println(c)
}
