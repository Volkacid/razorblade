package service

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateShortString(origURL string) string {
	hashGen := sha256.New()
	hashGen.Write([]byte(origURL))
	result := hex.EncodeToString(hashGen.Sum(nil))
	return result[:7]
}
