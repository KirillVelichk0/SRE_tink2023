package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	q := CreateCollector(200)
	config_data, err := ConstructConfigDataFromPath("../configs/config.yaml")
	if err != nil {
		log.Panicln(err.Error())
	}
	promCollector, err := ConstructPromCollector(logger, q)
	if err != nil {
		log.Panicln(err.Error())
	}
	host_url := ":" +
		strconv.FormatUint(uint64(config_data.ExporterPort), 10)
	var bind string
	flag.StringVar(&bind, "bind", host_url, "bind")
	flag.Parse()
	logger.Println(host_url)
	target_url := "http://" + config_data.TargetHost + ":" +
		strconv.FormatUint(uint64(config_data.TargetPort), 10)
	logger.Println(target_url)
	prometheus.Register(promCollector)
	prometheus.Register(promCollector.vecReqTime)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ctx := ctx
		for true {
			select {
			case <-ctx.Done():
				return
			default:
				logger.Printf("Going from executor to %s with count %d and limit %d", target_url, q.q.Length(), q.limit)
				if q.CanCollect() {
					stat := VizitApiWithGetRequest(target_url+"/api/v0/teams", logger)
					q.Collect([]RequestStatictic{stat})
					time.Sleep(time.Millisecond * 100)
				} else {
					time.Sleep(time.Second)
				}
			}

		}
	}()
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.HandlerFor(prometheus.Gatherers{
			prometheus.DefaultGatherer,
		}, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})
	logger.Println("Starting server")
	if err := http.ListenAndServe(bind, nil); err != nil {
		logger.Printf("Failed to start http server: %s", err)
	}
	cancel()
}
