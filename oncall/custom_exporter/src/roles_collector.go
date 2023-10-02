package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type DutyCollectorConfig struct {
	Labels     []string `yaml:"labels"`
	MetricName string   `yaml:"metricName"`
}

type DutyCollector struct {
	target_url string
	DutyCollectorConfig
	desk *prometheus.Desc
}

func ConstructDutyCollectorFromYaml(path string, target_url string) (*DutyCollector, error) {
	result := new(DutyCollector)
	filename, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config DutyCollectorConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	result.target_url = target_url
	result.DutyCollectorConfig = config
	return result, nil
}

// Describe prometheus describe
func (e *DutyCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Infof("Creating describe")
	e.desk = prometheus.NewDesc(e.MetricName, "This is roles info", e.Labels, nil)
}

// Collect prometheus collect
func (e *DutyCollector) Collect(ch chan<- prometheus.Metric) {
	log.Infof("Getting team list")
	teams, err := GetTeamsList(e.target_url)
	if err != nil {
		log.Errorf("Teams from url %s dont getted with error %s", e.target_url, err.Error())
		return
	}
	rangeGetter := new(CurTimeGetter)
	for _, team_name := range teams {
		team_roles, err := GetEventsFromTeamName(e.target_url, team_name, rangeGetter)
		if err != nil {
			log.Errorf("Roles from url %s and team %s dont getted with error %s",
				e.target_url, team_name, err.Error())
			return
		}
		for role, count := range team_roles {
			labelVals := []string{team_name, role}
			ch <- prometheus.MustNewConstMetric(e.desk, prometheus.GaugeValue, float64(count), labelVals...)
		}
	}

}
