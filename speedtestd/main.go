package main

import (
	_ "context"
	"flag"
	"time"

	"github.com/jrhorner1/ookla-speedtest/pkg/speedtest"
	log "github.com/sirupsen/logrus"
)

func ParseFlags() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yaml", "path to config")
	flag.Parse()
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}
	return configPath, nil
}

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

	switch config.Logging.Level {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	for {
		var results *speedtest.Speedtest
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
