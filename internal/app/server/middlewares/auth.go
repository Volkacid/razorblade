package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/google/uuid"
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
			stringValue := userID.Value[:len(userID.Value)-36] //without UUID
			userValue, _ := hex.DecodeString(stringValue)
			if !hmac.Equal(sign, userValue) {
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
	id := uuid.New()
	val := hex.EncodeToString(sign)
	val += id.String()
	return &http.Cookie{Name: "UserID", Value: val}
}

func createSign() []byte {
	hash := hmac.New(sha256.New, []byte(config.SecretKey))
	hash.Write([]byte("UserIdentificator"))
	return hash.Sum(nil)
}
