package random

import (
	"math/rand"
	"time"
)

func NewRandomString(n int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rnd.Intn(len(letter))]
	}
	return string(b)
}
