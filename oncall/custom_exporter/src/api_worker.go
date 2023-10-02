package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"time"
)

type TimeRangeGetter interface {
	GetTimeRange() (uint64, uint64)
}

type CurTimeGetter struct{}

func (time_getter *CurTimeGetter) GetTimeRange() (uint64, uint64) {
	curUnixTime := time.Now().Unix()
	return uint64(curUnixTime), uint64(curUnixTime)
}

func GetTeamsList(host_url string) ([]string, error) {
	resp, err := http.Get(host_url + "/api/v0/teams")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var json_data []string
	err = json.Unmarshal(body, &json_data)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, name := range json_data {
		result = append(result, name)
	}
	return result, nil
}

func GetEventsFromTeamName(host_url string, teamName string, timeRangeGetter TimeRangeGetter) (map[string]uint32, error) {
	if timeRangeGetter == nil {
		err := errors.New("empty interface")
		return nil, err
	}
	start_time, end_time := timeRangeGetter.GetTimeRange()
	start_time_str := strconv.FormatUint(start_time, 10)
	end_time_str := strconv.FormatUint(end_time, 10)
	counter := map[string]uint32{}
	counter["primary"] = 0
	counter["secondary"] = 0
	uri := host_url + "/api/v0/events?team=" + teamName + "&start__le=" + start_time_str +
		"&end__ge=" + end_time_str
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var json_data []map[string]interface{}
	err = json.Unmarshal(body, &json_data)
	if err != nil {
		return nil, err
	}
	for _, data := range json_data {
		role, ok := data["role"].(string)
		if !ok {
			err = &json.UnsupportedTypeError{}
			return nil, err
		}
		_, ok = counter[role]
		if ok {
			counter[role]++
		} else {
			counter[role] = 1
		}
	}

	return counter, nil
}
