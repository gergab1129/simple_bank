package util

import (
	"math/rand"
	"strings"
)

const abc = "abcdefghijklmnopqrstuxyz"

func randomInt(min, max int64) int64 {
	return min + max - rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var s strings.Builder

	for i := 0; i <= n; i++ {
		s.WriteString(string(abc[randomInt(0, 23)]))
	}

	return s.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return randomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{
		"EUR", "USD", "CAD", "COP",
	}

	return currencies[randomInt(0, 3)]
}
