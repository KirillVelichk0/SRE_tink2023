package main

import (
	"log"
	"net/http"
	"time"
)

func VizitApiWithGetRequest(url string, logger *log.Logger) RequestStatictic {
	var result RequestStatictic
	result.url = url
	startTime := time.Now()
	logger.Println("Go to api in func" + url)
	response, err := http.Get(url)
	if err != nil {
		result.success = false
		logger.Printf("Error while go to url %s\nError text: %s", url, err.Error())
		return result
	}
	defer response.Body.Close()

	responseTime := time.Since(startTime).Milliseconds()
	result.success = true
	result.statusCode = response.StatusCode
	result.responceTime = time.Duration(responseTime * int64(time.Millisecond))
	return result
}
