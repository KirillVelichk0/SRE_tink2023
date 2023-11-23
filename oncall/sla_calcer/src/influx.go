package main

import (
	"context"
	"errors"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxClient struct {
	c influxdb2.Client
}

func (c *InfluxClient) WriteToDb(name string, tags map[string]string, fields map[string]interface{}, org string, bucket string) error {
	writeApi := c.c.WriteAPIBlocking(org, bucket)
	p := influxdb2.NewPoint(name, tags, fields, time.Now())
	writeApi.WritePoint(context.Background(), p)
	return nil
}

func ConnectToInfluxDB() (influxdb2.Client, error) {

	dbToken := os.Getenv("INFLUXDB_TOKEN")
	if dbToken == "" {
		return nil, errors.New("INFLUXDB_TOKEN must be set")
	}

	dbURL := os.Getenv("INFLUXDB_URL")
	if dbURL == "" {
		return nil, errors.New("INFLUXDB_URL must be set")
	}

	client := influxdb2.NewClient(dbURL, dbToken)

	// validate client connection health
	_, err := client.Health(context.Background())

	return client, err
}

func CreateInfluxClient() (*InfluxClient, error) {
	var result *InfluxClient
	client, err := ConnectToInfluxDB()
	if err != nil {
		return result, err
	}
	result = new(InfluxClient)
	result.c = client
	return result, nil
}
