package main

import (
	"io"
	"math/rand"
	"net"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var charsetSize = len(charset)

func RandBytes(buf []byte) {
	for i := range buf {
		buf[i] = charset[rand.Intn(charsetSize)]
	}
}

func Swap(c1, c2 net.Conn) error {
	go func() {
		_, _ = io.Copy(c2, c1)
	}()
	_, err := io.Copy(c1, c2)
	return err
}
