// main.go

package main

import (
	"context"
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
		if len(os.Args) == 2 {
			cfg, _ := goconfig.LoadConfigFile(os.Args[1])
			fmt.Println("use conf")
			confPort, _ := cfg.GetValue("", "listenPort")
			listenPort = ":" + confPort
		} else {
			fmt.Println("use dafault conf")
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

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Println(s.Shutdown(context.TODO()))
}

//write request header to reponse header
func responseHeader(w http.ResponseWriter, r *http.Request) {
	for i, m := range r.Header {
		reposeVaule, _ := json.Marshal(m)
		w.Header().Add(i, string(reposeVaule))
	}
	w.Header().Add("VERSION", os.Getenv("VERSION"))
}
