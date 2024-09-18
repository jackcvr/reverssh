package main

import (
	"math/rand"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var charsetSize = len(charset)

func RandBytes(buf []byte) {
	for i := range buf {
		buf[i] = charset[rand.Intn(charsetSize)]
	}
}
