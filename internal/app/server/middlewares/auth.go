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
		userID, err := request.Cookie("UserID")
		if err != nil {
			createdCookie := createCookie(createSign())
			userID = createdCookie
			http.SetCookie(writer, createdCookie)
		} else {
			sign := createSign()
			userValue, _ := hex.DecodeString(userID.Value)
			if !hmac.Equal(sign, userValue[:len(userValue)-10]) {
				newCookie := createCookie(sign)
				userID = newCookie
				http.SetCookie(writer, newCookie)
			}
		}
		ctx := context.WithValue(request.Context(), config.UserID{}, userID.Value)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})

}

func createCookie(sign []byte) *http.Cookie {
	salt := make([]byte, 10)
	rand.Read(salt)
	sign = append(sign, salt...)
	return &http.Cookie{Name: "UserID", Value: hex.EncodeToString(sign)}
}

func createSign() []byte {
	hash := hmac.New(sha256.New, []byte(config.SecretKey))
	hash.Write([]byte("UserIdentificator"))
	return hash.Sum(nil)
}
