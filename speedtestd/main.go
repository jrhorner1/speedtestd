package main

import (
	_ "context"
	"flag"

	"github.com/jrhorner1/ookla-speedtest/pkg/speedtest"
	log "github.com/sirupsen/logrus"
)

func main() {
	var configPath string
	var retries int
	flag.StringVar(&configPath, "c", "./config.yaml", "path to config")
	flag.IntVar(&retries, "r", 3, "number of retries")
	flag.Parse()
	config, err := NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(logLevel(config.Logging.Level))

	var results *speedtest.Results
	if config.Speedtest.Server.Id != 0 {
		log.Debug("Running with server id")
		results = speedtest.RunWithServerId(config.Speedtest.Server.Id, retries)
	} else if config.Speedtest.Server.Name != "" {
		log.Debug("Running with server hostname")
		results = speedtest.RunWithHost(config.Speedtest.Server.Name, retries)
	} else {
		log.Debug("Running with default settings")
		results = speedtest.Run(retries)
	}
	influxdbConnect(results, config)
}

func logLevel(level string) log.Level {
	switch level {
	case "panic":
		return 0
	case "fatal":
		return 1
	case "error":
		return 2
	case "warn":
		return 3
	case "debug":
		return 5
	case "trace":
		return 6
	default: // info
		return 4
	}
}
