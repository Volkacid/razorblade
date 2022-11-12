package server

import "net/http"

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
