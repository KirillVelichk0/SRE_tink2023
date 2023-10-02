package main

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ConfigData struct {
	ExporterPort uint32 `yaml:"exporter_port"`
	ExporterHost string `yaml:"exporter_host"`
	TargetHost   string `yaml:"target_host"`
	TargetPort   uint32 `yaml:"target_port"`
}

func ConstructConfigDataFromPath(path string) (*ConfigData, error) {
	filename, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	configData := new(ConfigData)
	err = yaml.Unmarshal(yamlFile, configData)
	if err != nil {
		return nil, err
	}
	return configData, nil
}
