package util

import (
	"math/rand"
	"strings"
)

const symbols = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomVal(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomString(n int) string {
	b := strings.Builder{}
	b.Grow(n)

	for _ = range n {
		b.WriteByte(symbols[rand.Intn(len(symbols))])
	}

	return b.String()
}

func RandomOwner() string {
	return RandomString(10)
}

func RandomMoney() float64 {
	return RandomVal(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CNY", "RUB"}
	return currencies[rand.Intn(len(currencies))]
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
