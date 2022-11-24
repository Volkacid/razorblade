package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Volkacid/razorblade/internal/app/service"
	"net/http"
)

var secretKey = []byte("practicum")

func GetUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userIP := request.RemoteAddr
		//userIP, _, _ = strings.Cut(userIP, ":")
		sign := createSign(userIP)
		userID, err := request.Cookie("UserID")
		if err != nil {
			http.SetCookie(writer, createCookie(sign))
		} else if !hmac.Equal(sign, []byte(userID.Value)) {
			http.SetCookie(writer, createCookie(sign))
		}
		ctx := context.WithValue(request.Context(), service.UserID{}, hex.EncodeToString(sign))
		next.ServeHTTP(writer, request.WithContext(ctx))
	})

}

func createCookie(sign []byte) *http.Cookie {
	return &http.Cookie{Name: "UserID", Value: hex.EncodeToString(sign)}
}

func createSign(userIP string) []byte {
	hash := hmac.New(sha256.New, secretKey)
	hash.Write([]byte(userIP))
	return hash.Sum(nil)
}
