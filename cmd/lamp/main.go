package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	deviceIdAgentPool = map[string]*Agent{}
	systemIdAgentPool = map[string]*Agent{}
	initLog()
	go serve()
	go byteServe()
	monitor()
}

func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	logFile, err := os.Create("lamp.log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}
