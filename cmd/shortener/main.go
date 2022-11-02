package main

import (
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var db map[string]string
var mutex = &sync.RWMutex{}

func main() {
	db = make(map[string]string)
	rand.Seed(time.Now().Unix())

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Get("/", MainPage)
		router.Get("/{key}", GetHandler(db))
		router.Post("/", PostHandler(db))
	})
	http.ListenAndServe("localhost:8080", router)
}

func GetHandler(storage map[string]string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		key := chi.URLParam(request, "key")
		mutex.RLock()
		if storage[key] != "" {
			writer.Header().Set("Location", storage[key])
			writer.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			http.Error(writer, "Not found", http.StatusNotFound)
		}
		mutex.RUnlock()
	}
}

func PostHandler(storage map[string]string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var str string
		body, err := io.ReadAll(request.Body)
		defer request.Body.Close()
		if err != nil {
			http.Error(writer, "Please provide a correct link!", http.StatusBadRequest)
			return
		}

		str = string(body)
		if strings.Contains(str, "URL=") {
			_, str, _ = strings.Cut(string(body), "URL=")
		}
		str, _ = url.QueryUnescape(str)
		if !ValidateURL(str) {
			http.Error(writer, "Incorrect URL!", http.StatusBadRequest)
			return
		}
		mutex.Lock()
		for { //На случай, если сгенерированная последовательность уже будет занята
			foundStr := GenerateReducedLink()
			if storage[foundStr] == "" {
				storage[foundStr] = str
				writer.WriteHeader(http.StatusCreated)
				writer.Write([]byte("http://" + request.Host + "/" + foundStr))
				break
			}
		}
		mutex.Unlock()
	}
}

func GenerateReducedLink() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
	var builder strings.Builder
	for i := 0; i < 6; i++ {
		builder.WriteRune(chars[rand.Intn(len(chars))])
	}
	return builder.String()
}

func ValidateURL(str string) bool {
	path, err := url.ParseRequestURI(str)
	if err == nil && strings.ContainsAny(path.Host, ".:") {
		return true
	}
	return false
}

func MainPage(writer http.ResponseWriter, request *http.Request) {
	var form = `<html>
    <head>
    <title></title>
    </head>
    <body>
        <form name="shortener" method="post">
            <label>URL to short</label><input type="text" name="URL">
            <input type="submit" value="Shorten">
        </form>
    </body>
</html>`
	writer.Write([]byte(form))
}
