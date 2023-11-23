package main

import (
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.design/x/lockfree"
)

type RequestStatictic struct {
	statusCode   int
	success      bool
	responceTime time.Duration
	url          string
}

type RequestsResultsCollector struct {
	q     lockfree.Queue
	limit uint64
}

func CreateCollector(limit uint64) *RequestsResultsCollector {
	c := &RequestsResultsCollector{q: *lockfree.NewQueue(), limit: limit}
	return c
}

func (c *RequestsResultsCollector) CanCollect() bool {
	return c.q.Length() < c.limit
}

func (c *RequestsResultsCollector) Collect(s []RequestStatictic) error {
	if s == nil {
		return errors.New("s cant be nil")
	}
	for _, statistic := range s {
		c.q.Enqueue(statistic)
	}
	return nil
}

func (c *RequestsResultsCollector) TryGet() (RequestStatictic, error) {
	var result RequestStatictic
	val := c.q.Dequeue()
	if val == nil {
		return result, errors.New("queue is empty")
	}
	result, isOk := val.(RequestStatictic)
	if !isOk {
		return result, errors.New("Uncorrect contained type")
	}
	return result, nil
}

type PromCollector struct {
	logger            *log.Logger
	requestsCollector *RequestsResultsCollector
	vecReqTime        *prometheus.HistogramVec
	counter           uint64
}

func ConstructPromCollector(logger *log.Logger, reqCollector *RequestsResultsCollector) (*PromCollector, error) {
	var result *PromCollector
	if logger == nil || reqCollector == nil {
		return result, errors.New("Nil values")
	}
	result = new(PromCollector)
	result.counter = 0
	reqTimeOpts := prometheus.HistogramOpts{Name: "get_teams_resp_time", Namespace: "Oncall", Subsystem: "Oncall", Help: "Oncall get teams responce time sum count", Buckets: []float64{100, 500, 1000, 5000}}
	result.vecReqTime = prometheus.NewHistogramVec(reqTimeOpts, []string{"id", "statusCode"})
	result.logger = logger
	result.requestsCollector = reqCollector
	return result, nil
}

func (c *PromCollector) Describe(ch chan<- *prometheus.Desc) {
	c.logger.Println("Creating describer")
}

// Collect prometheus collect
func (c *PromCollector) Collect(ch chan<- prometheus.Metric) {
	c.logger.Println("Getting info")
	val := atomic.AddUint64(&c.counter, 1)
	//var stats []RequestStatictic
	for statistic, err := c.requestsCollector.TryGet(); err == nil; statistic, err = c.requestsCollector.TryGet() {
		if statistic.success {
			c.vecReqTime.WithLabelValues(fmt.Sprintf("%d", val), fmt.Sprintf("%d", statistic.statusCode)).Observe(float64(statistic.responceTime.Milliseconds()))
		} else {
			c.vecReqTime.WithLabelValues(fmt.Sprintf("%d", val), "_").Observe(1000000)
		}

	}
}
