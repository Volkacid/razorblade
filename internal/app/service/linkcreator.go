package service

import (
	"github.com/Volkacid/razorblade/internal/app/storage"
	"math/rand"
	"strings"
)

var chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")

func SetCreatorSeed(seed int64) {
	rand.Seed(seed)
}

func GenerateShortString(storage *storage.Storage) string {
	var builder strings.Builder
	for { //На случай, если сгенерированная последовательность уже будет занята
		for i := 0; i < 6; i++ {
			builder.WriteRune(chars[rand.Intn(len(chars))])
		}
		if _, err := storage.GetValue(builder.String()); err != nil {
			return builder.String()
		}
		builder.Reset()
	}
}
