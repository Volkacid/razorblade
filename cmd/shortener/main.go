package main

import (
	"bufio"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

//var storage map[string]string

func main() {
	//storage = make(map[string]string)
	http.HandleFunc("/", URLHandler)
	http.ListenAndServe("localhost:8080", nil)
}

func URLHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		query := request.URL.Path
		_, query, _ = strings.Cut(query, "/") //Обрезаем хост
		if query == "" {
			writer.Write([]byte(form))
			return
		}
		str := FindString(query)
		//str := storage[query]
		if str != "" {
			str, _, _ = strings.Cut(str, ":-:")
			writer.Header().Set("Location", str)
			writer.WriteHeader(http.StatusTemporaryRedirect)
			//writer.Write([]byte(str))
		} else {
			http.Error(writer, "Not found", http.StatusNotFound)
		}
		return
	case http.MethodPost:
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "Please provide a correct link!", http.StatusBadRequest)
			return
		}
		str := string(body)
		str, _ = url.QueryUnescape(str)
		if strings.Contains(str, "URL=") { //Если запрос из формы то обрезаем префикс
			_, str, _ = strings.Cut(str, "URL=")
		}
		foundStr := FindString(str)
		//foundStr := storage[str]
		if foundStr == "" {
			str = ReduceLink(str)
			writer.WriteHeader(http.StatusCreated)
			writer.Write([]byte("http://" + request.Host + "/" + str))
		} else {
			_, foundStr, _ = strings.Cut(foundStr, ":-:")
			writer.WriteHeader(http.StatusCreated)
			writer.Write([]byte("http://" + request.Host + "/" + foundStr))
		}
		return
	}
}

func FindString(message string) string {
	db, err := os.OpenFile("cmd/shortener/db.txt", os.O_RDONLY, 0444)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	scanner := bufio.NewScanner(db)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), message) {
			return scanner.Text()
		}
	}
	return ""
}

func ReduceLink(link string) string {
	rand.Seed(time.Now().Unix())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
	var builder strings.Builder
	for { //На случай, если сгенерированная последовательность уже будет занята
		for i := 0; i < 6; i++ {
			builder.WriteRune(chars[rand.Intn(len(chars))])
		}
		/*if storage[builder.String()] == "" {
			storage[builder.String()] = link
			break
		}*/
		if FindString(builder.String()) == "" {
			db, err := os.OpenFile("cmd/shortener/db.txt", os.O_WRONLY|os.O_APPEND, 0777)
			if err != nil {
				panic(err)
			}
			db.WriteString(link + ":-:" + builder.String() + "\n")
			defer db.Close()
			break
		}
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
