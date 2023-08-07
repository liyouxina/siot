package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

var Config *AllConfig

type MonitorConfig struct {
	Ip   string `json:"ip" yaml:"ip"`
	Port int    `json:"port" yaml:"port"`
}

type MingKongConfig struct {
	Port int `json:"port" yaml:"port"`
}

type AllConfig struct {
	MonitorConfig
	MingKongConfig
}

func init() {
	content, err := os.ReadFile("./config/mingkong.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(content, Config)
	if err != nil {
		panic(err)
	}
}
