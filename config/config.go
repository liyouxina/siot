package config

import (
	"encoding/json"
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

type LampConfig struct {
	Port int `json:"port" yaml:"port"`
}

type AllConfig struct {
	MonitorConfig  `json:"monitor_config"`
	MingKongConfig `json:"ming_kong_config"`
	LampConfig     `json:"lamp_config"`
}

func init() {
	content, err := os.ReadFile("./config/config.json")
	if err != nil {
		panic(err)
	}
	Config = &AllConfig{}
	err = json.Unmarshal(content, Config)
	if err != nil {
		panic(err)
	}
}
