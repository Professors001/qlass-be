package utils

import (
	"math/rand"
	"time"
)

const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // ตัดตัวอักษรที่สับสนง่ายออก (I, O, 0, 1)

func GenerateRandomString(length int) string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}