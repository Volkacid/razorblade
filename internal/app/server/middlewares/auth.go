package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Volkacid/razorblade/internal/app/config"
	"math/rand"
	"net/http"
)

func GetUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userIP := request.RemoteAddr
		//userIP, _, _ = strings.Cut(userIP, ":")
		sign := createSign(userIP)
		userID, err := request.Cookie("UserID")
		if err != nil {
			http.SetCookie(writer, createCookie(sign))
		} else {
			userValue, _ := hex.DecodeString(userID.Value)
			if !hmac.Equal(sign, userValue[:len(userValue)-5]) {
				http.SetCookie(writer, createCookie(sign))
			}
		}
		ctx := context.WithValue(request.Context(), config.UserID{}, hex.EncodeToString(sign))
		next.ServeHTTP(writer, request.WithContext(ctx))
	})

}

func createCookie(sign []byte) *http.Cookie {
	bytes := make([]byte, 5)
	rand.Read(bytes)
	sign = append(sign, bytes...)
	return &http.Cookie{Name: "UserID", Value: hex.EncodeToString(sign)}
}

func createSign(userIP string) []byte {
	hash := hmac.New(sha256.New, []byte(config.SecretKey))
	hash.Write([]byte(userIP))
	return hash.Sum(nil)
}
