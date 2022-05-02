package main

import (
	_ "context"
	"flag"
	"time"

	"github.com/jrhorner1/ookla-speedtest/pkg/speedtest"
	log "github.com/sirupsen/logrus"
)

func main() {
	configPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	config, err := NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config = envConfig(config)

	log.SetLevel(logLevel(config.Logging.Level))

	for {
		var results *speedtest.Results
		if config.Speedtest.Server.Id != 0 {
			log.Debug("Running with server id")
			results = speedtest.RunWithServerId(config.Speedtest.Server.Id)
		} else if config.Speedtest.Server.Name != "" {
			log.Debug("Running with server hostname")
			results = speedtest.RunWithHost(config.Speedtest.Server.Name)
		} else {
			log.Debug("Running with default settings")
			results = speedtest.Run()
		}
		influxdbConnect(results, config)
		log.Info("Sleeping for " + config.Speedtest.Interval + "...")
		intervalDuration, err := time.ParseDuration(config.Speedtest.Interval)
		if err != nil {
			log.Error("Sleep interval parse error:", err.Error())
		}
		log.Debug("Sleep Duration: ", intervalDuration)
		time.Sleep(intervalDuration)
	}
}

func ParseFlags() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yaml", "path to config")
	flag.Parse()
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}
	return configPath, nil
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
