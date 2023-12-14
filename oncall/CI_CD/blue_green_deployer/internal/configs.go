package internal

import (
	"os"

	"gopkg.in/yaml.v2"
)

type BlueGreenConfig struct {
	GreenUrl          string `yaml:"GreenUrl"`
	BlueUrl           string `yaml:"BlueUrl"`
	StartStateIsGreen bool   `yaml:"StartStateIsGreen"`
}
type ServerConfig struct {
	Port string `yaml:"Port"`
}

func ConstrServerCfgFromFile() (ServerConfig, error) {
	var result ServerConfig
	yamlFile, err := os.ReadFile("../configs/Server.yaml")
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func ConstrBlueGreenCfgFromFile() (BlueGreenConfig, error) {
	var result BlueGreenConfig
	yamlFile, err := os.ReadFile("../configs/BlueGreen.yaml")
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
