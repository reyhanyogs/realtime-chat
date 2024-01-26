package util

import (
	"math/rand"
	"strings"
)

func RandomString(n int) string {
	var sb strings.Builder
	alphabet := "abcdefghijklmopqrstuvxyz"
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
