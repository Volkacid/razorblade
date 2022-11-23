package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

var secretKey = []byte("practicum")

func GetUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userIP := request.RemoteAddr
		userIP, _, _ = strings.Cut(userIP, ":")
		userID, _ := request.Cookie("UserID")

		if userID == nil || !validateUser(userIP, userID.Value) {
			sign := createSign(userIP)
			cookie := &http.Cookie{Name: "UserID", Value: hex.EncodeToString(sign)}
			http.SetCookie(writer, cookie)
		}
		next.ServeHTTP(writer, request)
	})

}

func createSign(userIP string) []byte {
	hash := hmac.New(sha256.New, secretKey)
	hash.Write([]byte(userIP))
	return hash.Sum(nil)
}

func validateUser(userIP string, userID string) bool {
	sign := createSign(userIP)
	return hmac.Equal(sign, []byte(userID))
}
