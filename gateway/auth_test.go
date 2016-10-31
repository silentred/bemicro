package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	sign := makeSign()
	assert.True(t, isValid(sign))
}

func BenchmarkSign(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sign := makeSign()
		//isValid(sign)
		_ = sign
	}
}
