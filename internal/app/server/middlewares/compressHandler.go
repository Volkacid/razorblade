package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (gzWriter gzipWriter) Write(bytes []byte) (int, error) {
	return gzWriter.Writer.Write(bytes)
}

func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		newWriter := writer
		newRequest := request

		if strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			gzWriter, err := gzip.NewWriterLevel(writer, gzip.BestSpeed)
			if err != nil {
				http.Error(writer, "Compression error", http.StatusInternalServerError)
				return
			}
			defer gzWriter.Close()
			writer.Header().Set("Content-Encoding", "gzip")
			newWriter = gzipWriter{ResponseWriter: newWriter, Writer: gzWriter}
		}

		if strings.Contains(request.Header.Get("Content-Encoding"), "gzip") {
			gzReader, err := gzip.NewReader(request.Body)
			if err != nil {
				http.Error(writer, "Decompression error", http.StatusInternalServerError)
				return
			}
			defer gzReader.Close()
			newRequest.Body = gzReader
		}

		next.ServeHTTP(newWriter, newRequest)
	})
}
