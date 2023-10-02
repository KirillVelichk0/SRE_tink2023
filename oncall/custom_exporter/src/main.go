package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {

	config_data, err := ConstructConfigDataFromPath("../configs/service_config.yaml")
	if err != nil {
		log.Errorf("Cant parse data from service_config.yaml with error %s", err.Error())
		os.Exit(1)
	}
	host_url := ":" +
		strconv.FormatUint(uint64(config_data.ExporterPort), 10)
	log.Infof(host_url)
	var bind string
	flag.StringVar(&bind, "bind", host_url, "bind")
	flag.Parse()

	target_url := "http://" + config_data.TargetHost + ":" +
		strconv.FormatUint(uint64(config_data.TargetPort), 10)
	log.Infof(target_url)
	collector, err := ConstructDutyCollectorFromYaml("../configs/metric_struct.yaml", target_url)
	if err != nil {
		log.Errorf("Cant parse data from metric_struct.yaml with error %s", err.Error())
		os.Exit(1)
	}
	prometheus.Register(collector)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.HandlerFor(prometheus.Gatherers{
			prometheus.DefaultGatherer,
		}, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})

	log.Infof("Starting http server - %s", bind)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Errorf("Failed to start http server: %s", err)
	}
	log.Infof("server started at %s", bind)
}
