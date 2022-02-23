// main.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Unknwon/goconfig"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		responseHeader(rw, r)
		fmt.Fprintf(rw, "Hello World, %v\n", time.Now())
	})
	http.HandleFunc("/healthz", func(rw http.ResponseWriter, r *http.Request) {
		responseHeader(rw, r)
		fmt.Fprint(rw, "server is healthy")
	})

	listenPort := os.Getenv("listenPort")
	if listenPort != "" {
		fmt.Print("user env listenPort")
		listenPort = ":" + listenPort
	} else {
		cfg, err := goconfig.LoadConfigFile("http.conf")
		if err == nil {
			fmt.Println("use conf")
			confPort, _ := cfg.GetValue("", "listenPort")
			listenPort = ":" + confPort
		} else {
			fmt.Print("use dafault 80")
			listenPort = ":80"
		}
	}

	fmt.Println(listenPort)
	s := &http.Server{
		Addr:           listenPort,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		log.Println(s.ListenAndServe())
		log.Println("server shutdown")
	}()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	// Stop the service gracefully.
	log.Println(s.Shutdown(nil))
	// Wait gorotine print shutdown message
}

//write request header to reponse header
func responseHeader(w http.ResponseWriter, r *http.Request) {
	for i, m := range r.Header {
		reposeVaule, _ := json.Marshal(m)
		w.Header().Add(i, string(reposeVaule))
	}
	w.Header().Add("VERSION", os.Getenv("VERSION"))
}
