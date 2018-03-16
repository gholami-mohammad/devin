package helpers

import (
	"math/rand"
	"time"
)

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-")
	var str = make([]rune, length)

	for i := range str {
		rand.Seed(time.Now().UnixNano())
		str[i] = letters[rand.Intn(len(letters))]
	}

	return string(str)
}
