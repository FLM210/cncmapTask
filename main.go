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
	// delayTime := prometheus.NewHistogram(prometheus.HistogramOpts{
	// 	Name:    "http_request_durations",
	// 	Help:    "A histogram of the HTTP request durations in seconds.",
	// 	Buckets: prometheus.ExponentialBuckets(0.1, 2, 5),
	// })

	processingtime := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "our_company",
		Subsystem: "blob_storage",
		Name:      "http_request_processingtime",
		Help:      "Number of blob storage operations waiting to be processed.",
	})

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		time.Sleep(time.Duration(createRandomInt(2) * 1e9))
		responseHeader(rw, r)
		fmt.Fprintf(rw, "Hello World, %v\n", time.Now())
		costTime := time.Since(startTime)
		fmt.Println(costTime.Seconds())
		processingtime.Set(costTime.Seconds())
	})
	http.HandleFunc("/healthz", func(rw http.ResponseWriter, r *http.Request) {
		responseHeader(rw, r)
		fmt.Fprintf(rw, "Hello World, %v\n", time.Now())
		responseHeader(rw, r)
		fmt.Fprint(rw, "server is healthy")
	})
	prometheus.Register(opsQueued)
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
