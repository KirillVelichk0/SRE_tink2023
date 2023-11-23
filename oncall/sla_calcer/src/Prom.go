package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type PrometheusConfig struct {
	Host string `yaml:"prom_host"`
	Port int    `yaml:"prom_port"`
}

type Metric struct {
	Le string `json:"le"`
}

type TimeRespDataResultVal struct {
	MetricVal Metric        `json:"metric"`
	Value     []interface{} `json:"value"`
}

type TimeRespData struct {
	ResultType string                  `json:"resultType"`
	Result     []TimeRespDataResultVal `json:"result"`
}

type TimeResp struct {
	Status string       `json:"status"`
	Data   TimeRespData `json:"data"`
}

func FormatTimeResp(r TimeResp) (FormattedCollectionData, error) {
	var result FormattedCollectionData
	for _, data := range r.Data.Result {
		val, err := strconv.ParseUint(data.Value[1].(string), 10, 64)
		if err != nil {
			return result, err
		}
		dataF := FormattedData{unixT: data.Value[0].(float64),
			val: val}
		switch data.MetricVal.Le {
		case "100":
			result.f100 = dataF
		case "500":
			result.f500 = dataF
		case "1000":
			result.f1000 = dataF
		case "5000":
			result.f5000 = dataF
		default:
			result.fBig = dataF
		}
	}
	return result, nil
}

type FormattedData struct {
	unixT float64
	val   uint64
}
type FormattedCollectionData struct {
	f100  FormattedData
	f500  FormattedData
	f1000 FormattedData
	f5000 FormattedData
	fBig  FormattedData
}

func ParseConfig(path string) (PrometheusConfig, error) {
	var result PrometheusConfig
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(yamlFile, &result)
	return result, err
}

type PrometheusVisitor struct {
	cfg PrometheusConfig
}

func (v *PrometheusVisitor) GoToProm(query string) (int, []byte, error) {
	var result []byte
	curTime := time.Now().Unix()
	url := fmt.Sprintf("http://%s:%d/api/v1/query?query=%s&time=%d", v.cfg.Host, v.cfg.Port, query, curTime)
	resp, err := http.Get(url)
	if err != nil {
		return 0, result, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, result, err
	}

	return resp.StatusCode, data, nil
}

type PrometheusServicer struct {
	v PrometheusVisitor
}

func CreateServicer(cfg PrometheusConfig) *PrometheusServicer {
	result := new(PrometheusServicer)
	result.v.cfg = cfg
	return result
}

func (s *PrometheusServicer) GetTimeToResponce() (TimeResp, error) {
	var result TimeResp
	statusCode, data, err := s.v.GoToProm("sum+by+(le)(Oncall_Oncall_get_teams_resp_time_bucket)")
	if err != nil {
		return result, err
	}
	if statusCode != 200 {
		return result, fmt.Errorf("Status code is not 200 - %d", statusCode)
	}
	err = json.Unmarshal(data, &result)
	return result, err
}
