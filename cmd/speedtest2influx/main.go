package main

import (
	_ "context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/jrhorner1/ookla-speedtest/pkg/speedtest"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Influxdb struct {
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"influxdb"`
	Speedtest struct {
		Server struct {
			Id   int    `yaml:"id"`
			Name string `yaml:"name"`
		}
		Interval string `yaml:"interval"`
	} `yaml:"speedtest"`
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", path)
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		intvalue, _ := strconv.Atoi(value)
		return intvalue
	}
	return fallback
}

func envConfig(config *Config) *Config {
	newConfig := &Config{}
	newConfig.Influxdb.Address = getEnv("INFLUXDB_ADDRESS", config.Influxdb.Address)
	newConfig.Influxdb.Port = getEnvInt("INFLUXDB_PORT", config.Influxdb.Port)
	newConfig.Influxdb.Database = getEnv("INFLUXDB_DATABASE", config.Influxdb.Database)
	newConfig.Influxdb.Username = getEnv("INFLUXDB_USERNAME", config.Influxdb.Username)
	newConfig.Influxdb.Password = getEnv("INFLUXDB_PASSWORD", config.Influxdb.Password)
	newConfig.Speedtest.Server.Id = getEnvInt("SPEEDTEST_SERVER_ID", config.Speedtest.Server.Id)
	newConfig.Speedtest.Server.Name = getEnv("SPEEDTEST_SERVER_NAME", config.Speedtest.Server.Name)
	newConfig.Speedtest.Interval = getEnv("SPEEDTEST_INTERVAL", config.Speedtest.Interval)
	newConfig.Logging.Level = getEnv("LOGGING_LEVEL", config.Logging.Level)
	return newConfig
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

func influxdbConnect(results *speedtest.Speedtest, config *Config) {
	influxdb_protocol := "http"
	influxdb_server := config.Influxdb.Address
	influxdb_port := config.Influxdb.Port
	influxdb_url := influxdb_protocol + "://" + influxdb_server + ":" + strconv.Itoa(influxdb_port)
	influxdb_user := config.Influxdb.Username
	influxdb_pass := config.Influxdb.Password
	influxdb_token := influxdb_user + ":" + influxdb_pass
	influxdb_org := ""
	influxdb_database := config.Influxdb.Database

	log.Info("Connecting to influxdb server: " + influxdb_url)
	client := influxdb2.NewClient(influxdb_url, influxdb_token)

	writeAPI := client.WriteAPI(influxdb_org, influxdb_database)
	errorsCh := writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			log.Error("write error:", err.Error())
		}
	}()

	tags := map[string]string{
		"serverId":       strconv.Itoa(results.Server.Id),
		"serverName":     results.Server.Name,
		"serverLocation": results.Server.Location,
		"serverCountry":  results.Server.Country,
		"serverHost":     results.Server.Host,
		"serverPort":     strconv.Itoa(results.Server.Port),
		"serverIp":       results.Server.Ip,
		"isp":            results.Isp}

	ping := influxdb2.NewPoint(
		"ping", // measurement
		tags,
		map[string]interface{}{ // fields
			"jitter":  results.Ping.Jitter,
			"latency": results.Ping.Latency},
		results.Timestamp)
	log.Debug("Writing ping measurements")
	writeAPI.WritePoint(ping)

	download := influxdb2.NewPoint(
		"download", // measurement
		tags,
		map[string]interface{}{
			"bandwidth": results.Download.Bandwidth * 8, // Value is in bytes, converting to bits
			"bytes":     results.Download.Bytes,
			"elapsed":   results.Download.Elapsed},
		results.Timestamp)
	log.Debug("Writing download measurements")
	writeAPI.WritePoint(download)

	upload := influxdb2.NewPoint(
		"upload", // measurement
		tags,
		map[string]interface{}{
			"bandwidth": results.Upload.Bandwidth * 8, // Value is in bytes, converting to bits
			"bytes":     results.Upload.Bytes,
			"elapsed":   results.Upload.Elapsed},
		results.Timestamp)
	log.Debug("Writing upload measurements")
	writeAPI.WritePoint(upload)

	packet := influxdb2.NewPoint(
		"packet", // measurement
		tags,
		map[string]interface{}{ // fields
			"loss": results.PacketLoss},
		results.Timestamp)
	log.Debug("Writing packet measurements")
	writeAPI.WritePoint(packet)

	clientData := influxdb2.NewPoint(
		"client",
		tags,
		map[string]interface{}{
			"internalIp":       results.Interface.InternalIp,
			"interfaceName":    results.Interface.Name,
			"interfaceMacAddr": results.Interface.MacAddr,
			"isVpn":            strconv.FormatBool(results.Interface.IsVpn),
			"externalIp":       results.Interface.ExternalIp,
			"isp":              results.Isp},
		results.Timestamp)
	log.Debug("Writing client info")
	writeAPI.WritePoint(clientData)

	serverData := influxdb2.NewPoint(
		"server",
		tags,
		map[string]interface{}{
			"serverId":       strconv.Itoa(results.Server.Id),
			"serverName":     results.Server.Name,
			"serverLocation": results.Server.Location,
			"serverCountry":  results.Server.Country,
			"serverHost":     results.Server.Host,
			"serverPort":     strconv.Itoa(results.Server.Port),
			"serverIp":       results.Server.Ip},
		results.Timestamp)
	log.Debug("Writing server info")
	writeAPI.WritePoint(serverData)

	result := influxdb2.NewPoint(
		"result",
		tags,
		map[string]interface{}{
			"id":  results.Result.Id,
			"url": results.Result.Url},
		results.Timestamp)
	log.Debug("Writing result info")
	writeAPI.WritePoint(result)

	log.Info("Flushing writes from the buffer")
	writeAPI.Flush()
	log.Info("Closing the influxdb client connection")
	client.Close()
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
