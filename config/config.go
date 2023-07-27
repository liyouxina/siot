package config

import "os"

type Monitor struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type MingKongConfig struct {
	Port    int `json:"port"`
	Monitor `json:"monitor"`
}

func GetMingKongConfig() {
	content, err := os.ReadFile("./config/mingkong.yml")
	yaml
}
