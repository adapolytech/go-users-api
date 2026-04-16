package server

import (
	"io"
	"log"
	"net/http"
)

func CreateServer() *http.ServeMux {
	server := http.NewServeMux()
	return server
}

var Logger http.HandlerFunc = func(resp http.ResponseWriter, req *http.Request) {
	data, _ := io.ReadAll(req.Body)
	log.Default().Printf("URL: %s, Content: %s", req.URL.String(), string(data))
}

var LoggerMiddleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Default().Printf("URL: %s, Content: %s", req.URL.String(), req.Method)
		next.ServeHTTP(w, req)
	}
}
