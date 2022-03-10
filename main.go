// version: v3

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	delayTime := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_durations",
		Help:    "A histogram of the HTTP request durations in seconds.",
		Buckets: prometheus.ExponentialBuckets(0.1, 2, 5),
	})

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(delayTime)
		responseHeader(rw, r)
		fmt.Fprintf(rw, "Hello World, %v\n", time.Now())
		timer.ObserveDuration()
	})
	http.HandleFunc("/healthz", func(rw http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(delayTime)
		defer timer.ObserveDuration()
		responseHeader(rw, r)
		fmt.Fprintf(rw, "Hello World, %v\n", time.Now())
		responseHeader(rw, r)
		fmt.Fprint(rw, "server is healthy")
	})
	prometheus.Register(delayTime)
	http.Handle("/metrics", promhttp.Handler())

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

	//优雅终止
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

//generate random numbers
func createRandomInt(num int) int {
	r := rand.Intn(num)
	return r
}
