package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type zipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w zipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func ZipHandlerWrite(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(zipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func ZipHandlerRead(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()
		r.Body = gz
		next.ServeHTTP(w, r)
	})
}
