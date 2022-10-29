package main

import (
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var db map[string]string

func main() {
	db = make(map[string]string)
	http.HandleFunc("/", URLHandler(db))
	http.ListenAndServe("localhost:8080", nil)
}

func URLHandler(storage map[string]string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			query := request.URL.Path
			_, query, _ = strings.Cut(query, "/") //Обрезаем хост
			if query == "" {
				writer.Write([]byte(form))
				return
			}
			str := storage[query]
			if str != "" {
				writer.Header().Set("Location", str)
				writer.WriteHeader(http.StatusTemporaryRedirect)
			} else {
				http.Error(writer, "Not found", http.StatusNotFound)
			}
			return
		case http.MethodPost:
			var str string
			body, err := io.ReadAll(request.Body)
			if err != nil {
				http.Error(writer, "Please provide a correct link!", http.StatusBadRequest)
				return
			}
			str = string(body)
			if strings.Contains(str, "URL=") {
				_, str, _ = strings.Cut(string(body), "URL=")
			}
			str, _ = url.QueryUnescape(str)
			/*if path, _ := url.ParseRequestURI(str); path == nil {
				http.Error(writer, "Incorrect link! Make sure it begins with http:// or https://", http.StatusBadRequest)
				return
			}*/
			for { //На случай, если сгенерированная последовательность уже будет занята
				foundStr := GenerateReducedLink()
				if storage[foundStr] == "" {
					storage[foundStr] = str
					writer.WriteHeader(http.StatusCreated)
					writer.Write([]byte("http://" + request.Host + "/" + foundStr))
					break
				}
			}
			return
		}
	}
}

func GenerateReducedLink() string {
	rand.Seed(time.Now().Unix())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
	var builder strings.Builder
	for i := 0; i < 6; i++ {
		builder.WriteRune(chars[rand.Intn(len(chars))])
	}
	return builder.String()
}

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
