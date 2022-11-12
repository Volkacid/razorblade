package service

import (
	"math/rand"
	"strings"
)

var chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")

func SetCreatorSeed(seed int64) {
	rand.Seed(seed)
}

func GenerateShortString() string {
	var builder strings.Builder
	for i := 0; i < 6; i++ {
		builder.WriteRune(chars[rand.Intn(len(chars))])
	}
	return builder.String()
}
