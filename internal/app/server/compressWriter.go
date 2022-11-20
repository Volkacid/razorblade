package server

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
		if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(writer, request)
			return
		}

		gzWriter, err := gzip.NewWriterLevel(writer, gzip.BestSpeed)
		if err != nil {
			http.Error(writer, "Compression error", http.StatusInternalServerError)
			return
		}
		defer gzWriter.Close()
		writer.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: writer, Writer: gzWriter}, request)
	})
}
