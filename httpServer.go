package main

import (
	"encoding/json"
	"io"

	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	http.Handle("/", index(indexhandler()))
	http.Handle("/healthz", healthz(healthzhandler()))
	LISTENPORT := os.Getenv("LISTENPORT")
	if LISTENPORT != "" {
		http.ListenAndServe(":"+LISTENPORT, nil)
	} else {
		http.ListenAndServe(":80", nil)
	}
}

func indexhandler() http.Handler {
	return http.HandlerFunc(indexHandler)
}
func indexHandler(rw http.ResponseWriter, r *http.Request) {
	responseHeader(rw, r)
	io.WriteString(rw, "Hello World")
}
func index(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func healthzhandler() http.Handler {
	return http.HandlerFunc(healthzHandler)
}
func healthzHandler(rw http.ResponseWriter, r *http.Request) {
	responseHeader(rw, r)
	io.WriteString(rw, "server is ok")
}
func healthz(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

//write request header to reponse header
func responseHeader(w http.ResponseWriter, r *http.Request) {
	for i, m := range r.Header {
		reposeVaule, _ := json.Marshal(m)
		w.Header().Add(i, string(reposeVaule))
	}
	w.Header().Add("VERSION", os.Getenv("VERSION"))
}
